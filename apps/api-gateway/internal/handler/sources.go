package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// GetChatSources はチャットで使用されたソースを取得
func (h *Handler) GetChatSources(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	chatId openapi_types.UUID,
) {
	ctx := r.Context()

	// Step 1: チャットの存在確認
	_, err := h.queries.GetChat(ctx, db.GetChatParams{
		ID:          chatId,
		WorkspaceID: workspaceId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "CHAT_NOT_FOUND", "Chat not found")
			return
		}
		log.Printf("Failed to get chat: %v", err)
		respondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch chat")
		return
	}

	// Step 2: チャットメッセージからdocument_refsを収集
	messages, err := h.queries.GetChatMessages(ctx, db.GetChatMessagesParams{
		ChatID: chatId,
		Limit:  1000, // 全メッセージ取得
	})
	if err != nil {
		log.Printf("Failed to get messages: %v", err)
		respondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch messages")
		return
	}

	// Step 3: 全メッセージのdocument_refsを結合
	allRefs := []byte("[]")
	if len(messages) > 0 {
		// 複数メッセージのdocument_refsを統合
		var combinedRefs []json.RawMessage
		for _, msg := range messages {
			if msg.DocumentRefs.Valid && len(msg.DocumentRefs.RawMessage) > 0 {
				combinedRefs = append(combinedRefs, msg.DocumentRefs.RawMessage)
			}
		}

		// TODO: 複数のJSONB配列を1つに統合する処理
		// 今は最初のメッセージのみ使用（仮実装）
		if len(messages) > 0 && messages[0].DocumentRefs.Valid {
			allRefs = messages[0].DocumentRefs.RawMessage
		}
	}

	// Step 4: SourceServiceを呼び出し
	sources, err := h.sourceService.GetSources(ctx, workspaceId, allRefs)
	if err != nil {
		log.Printf("Failed to get sources: %v", err)
		respondError(w, http.StatusInternalServerError, "SOURCES_ERROR", "Failed to get sources")
		return
	}

	respondJSON(w, http.StatusOK, sources)
}

// GetAnalysisSources は分析で使用されたソースを取得
func (h *Handler) GetAnalysisSources(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	analysisId openapi_types.UUID,
) {
	// ctx := r.Context()

	// TODO: 分析の実装後に追加
	respondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Analysis sources not yet implemented")
}
