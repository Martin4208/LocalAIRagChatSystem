package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/google/uuid"
)

type SourceService struct {
	queries *db.Queries
}

func NewSourceService(queries *db.Queries) *SourceService {
	return &SourceService{queries: queries}
}

// DocumentReference はJSONBから解析する内部型
// page_number を追加（nullable: 旧データとの後方互換を保つ）
type DocumentReference struct {
	DocumentID     uuid.UUID `json:"document_id"`
	DocumentName   string    `json:"document_name"`
	ChunkIndex     int32     `json:"chunk_index"`
	PageNumber     *int      `json:"page_number,omitempty"` // ← 追加
	Score          float64   `json:"score"`
	ContentPreview string    `json:"content_preview"`
}

type DocumentMetadata struct {
	ID        uuid.UUID
	Name      string
	MimeType  string
	SizeBytes int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetSources はJSONBからソース情報を構築する
func (s *SourceService) GetSources(
	ctx context.Context,
	workspaceID uuid.UUID,
	documentRefsJSON []byte,
) (*api.SourcesResponse, error) {

	refs, err := s.parseDocumentReferences(documentRefsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document_refs: %w", err)
	}

	if len(refs) == 0 {
		return &api.SourcesResponse{
			Sources:        []api.SourceDocument{},
			TotalDocuments: 0,
			TotalChunks:    0,
		}, nil
	}

	uniqueDocIDs := s.extractUniqueDocumentIDs(refs)

	metadata, err := s.fetchDocumentMetadata(ctx, workspaceID, uniqueDocIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document metadata: %w", err)
	}

	sources := s.groupByDocument(refs, metadata)

	return &api.SourcesResponse{
		Sources:        sources,
		TotalDocuments: int32(len(sources)),
		TotalChunks:    int32(len(refs)),
	}, nil
}

// parseDocumentReferences はJSONBをパースする
func (s *SourceService) parseDocumentReferences(data []byte) ([]DocumentReference, error) {
	if len(data) == 0 || string(data) == "null" {
		return []DocumentReference{}, nil
	}

	var refs []DocumentReference
	if err := json.Unmarshal(data, &refs); err != nil {
		return nil, err
	}
	return refs, nil
}

// extractUniqueDocumentIDs はユニークなdocument_idを抽出
func (s *SourceService) extractUniqueDocumentIDs(refs []DocumentReference) []uuid.UUID {
	seen := make(map[uuid.UUID]bool)
	var uniqueIDs []uuid.UUID
	for _, ref := range refs {
		if !seen[ref.DocumentID] {
			seen[ref.DocumentID] = true
			uniqueIDs = append(uniqueIDs, ref.DocumentID)
		}
	}
	return uniqueIDs
}

// fetchDocumentMetadata はDBからドキュメント情報を取得
func (s *SourceService) fetchDocumentMetadata(
	ctx context.Context,
	workspaceID uuid.UUID,
	documentIDs []uuid.UUID,
) (map[uuid.UUID]*DocumentMetadata, error) {
	rows, err := s.queries.GetDocumentMetadataByIDs(ctx, db.GetDocumentMetadataByIDsParams{
		DocumentIds: documentIDs,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]*DocumentMetadata)
	for _, row := range rows {
		result[row.ID] = &DocumentMetadata{
			ID:        row.ID,
			Name:      row.Name,
			MimeType:  row.MimeType,
			SizeBytes: row.SizeBytes,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}
	return result, nil
}

// groupByDocument はドキュメントごとにチャンクをグルーピングし、
// referenced_pages（参照されたページ番号の重複排除リスト）も構築する。
func (s *SourceService) groupByDocument(
	refs []DocumentReference,
	metadata map[uuid.UUID]*DocumentMetadata,
) []api.SourceDocument {

	grouped := make(map[uuid.UUID][]api.ChunkReference)
	// ドキュメントごとにページ番号のセットを管理
	referencedPages := make(map[uuid.UUID]map[int]struct{})

	for _, ref := range refs {
		// page_number があればセットに追加
		if ref.PageNumber != nil {
			if referencedPages[ref.DocumentID] == nil {
				referencedPages[ref.DocumentID] = make(map[int]struct{})
			}
			referencedPages[ref.DocumentID][*ref.PageNumber] = struct{}{}
		}

		chunk := api.ChunkReference{
			ChunkIndex:     int32(ref.ChunkIndex),
			PageNumber:     ref.PageNumber, // ← 追加
			ContentPreview: ref.ContentPreview,
			RelevanceScore: &ref.Score,
		}
		grouped[ref.DocumentID] = append(grouped[ref.DocumentID], chunk)
	}

	var sources []api.SourceDocument

	for docID, chunks := range grouped {
		meta := metadata[docID]
		if meta == nil {
			continue
		}

		// チャンクをスコア降順でソート
		sort.Slice(chunks, func(i, j int) bool {
			if chunks[i].RelevanceScore == nil {
				return false
			}
			if chunks[j].RelevanceScore == nil {
				return true
			}
			return *chunks[i].RelevanceScore > *chunks[j].RelevanceScore
		})

		// referenced_pages をソート済みスライスに変換
		pages := s.sortedPages(referencedPages[docID])

		sources = append(sources, api.SourceDocument{
			DocumentId:      docID,
			DocumentName:    meta.Name,
			MimeType:        meta.MimeType,
			SizeBytes:       meta.SizeBytes,
			ReferencedPages: pages, // ← 追加：フロントエンドで「P.2, P.5 を参照」と表示できる
			ChunksUsed:      chunks,
			CreatedAt:       meta.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       meta.UpdatedAt.Format(time.RFC3339),
		})
	}

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].DocumentName < sources[j].DocumentName
	})

	return sources
}

// sortedPages はページ番号セットをソート済みスライスに変換するヘルパー
func (s *SourceService) sortedPages(pageSet map[int]struct{}) []int {
	if pageSet == nil {
		return nil
	}
	pages := make([]int, 0, len(pageSet))
	for p := range pageSet {
		pages = append(pages, p)
	}
	sort.Ints(pages)
	return pages
}
