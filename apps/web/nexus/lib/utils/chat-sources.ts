// lib/utils/chat-sources.ts

import { ChatMessage, SourceMapping } from '@/lib/types/chat';
import { SourceDocument, ChunkReference } from '@/lib/types/sources';

interface DocumentRefRaw {
  document_id: string;
  document_name?: string;
  chunk_index: number;
  score: number;
  content_preview?: string;
}

interface ExtractedSources {
  allSources: SourceDocument[];
  messageSourcesMap: Map<string, number[]>;
}


/**
 * 全メッセージから重複なしのソースリストを作成
 * document_id + chunk_index で一意性を判定
 * 最初に出現した順で番号を振る
 */
export function extractUniqueSources(messages: ChatMessage[]): SourceMapping[] {
  const sourceMap = new Map<string, SourceMapping>();
  let currentIndex = 1;

  // メッセージを古い順に処理（message_indexでソート）
  const sortedMessages = [...messages].sort((a, b) => a.message_index - b.message_index);

  for (const message of sortedMessages) {
    if (!message.documentRefs || message.documentRefs.length === 0) {
      continue;
    }

    for (const ref of message.documentRefs) {
      // ユニークキー: document_id + chunk_index
      const key = `${ref.document_id}-${ref.chunk_index}`;

      // まだ登録されていなければ追加
      if (!sourceMap.has(key)) {
        sourceMap.set(key, {
          index: currentIndex++,
          documentId: ref.document_id,
          documentName: ref.document_name || null,
          chunkIndex: ref.chunk_index,
          score: ref.score,
          contentPreview: ref.content_preview || null,
        });
      }
    }
  }

  // Map → Array
  return Array.from(sourceMap.values());
}

/**
 * 特定のメッセージが参照しているソース番号を取得
 */
export function mapMessageToSources(
  message: ChatMessage,
  sourcesList: SourceMapping[]
): number[] {
  if (!message.documentRefs || message.documentRefs.length === 0) {
    return [];
  }

  const indices: number[] = [];

  for (const ref of message.documentRefs) {
    const key = `${ref.document_id}-${ref.chunk_index}`;
    
    // sourcesListから該当するものを探す
    const source = sourcesList.find(
      (s) => `${s.documentId}-${s.chunkIndex}` === key
    );

    if (source) {
      indices.push(source.index);
    }
  }

  // 番号順にソート
  return indices.sort((a, b) => a - b);
}

/**
 * ソース番号から詳細情報を取得
 */
export function findSourceByIndex(
  index: number,
  sourcesList: SourceMapping[]
): SourceMapping | undefined {
  return sourcesList.find((s) => s.index === index);
}

/**
 * 全メッセージのソース番号マッピングを作成
 */
export function buildMessageSourcesMap(
  messages: ChatMessage[],
  sourcesList: SourceMapping[]
): Map<string, number[]> {
  const map = new Map<string, number[]>();

  for (const message of messages) {
    const indices = mapMessageToSources(message, sourcesList);
    if (indices.length > 0) {
      map.set(message.id, indices);
    }
  }

  return map;
}


/**
 * チャットメッセージから引用ソースを抽出して整形
 */
export function extractSourcesFromMessages(messages: ChatMessage[]): ExtractedSources {
  const allSources: SourceDocument[] = [];
  const messageSourcesMap = new Map<string, number[]>();
  
  // document_id ごとにグループ化するためのマップ
  const documentMap = new Map<string, {
    document: Partial<SourceDocument>;
    chunks: ChunkReference[];
    sourceIndex?: number;
  }>();

  // 全メッセージから documentRefs を収集
  messages.forEach((message) => {
    if (message.role !== 'assistant' || !message.documentRefs) {
      return;
    }

    const messageSourceIndices: number[] = [];

    message.documentRefs.forEach((ref) => {
      const docId = ref.document_id;

      // 既存のドキュメントエントリを取得または作成
      if (!documentMap.has(docId)) {
        documentMap.set(docId, {
          document: {
            document_id: docId,
            document_name: ref.document_name || 'Unknown Document',
            mime_type: 'application/pdf', // TODO: 実際のmime_typeを取得
            size_bytes: 0, // TODO: 実際のサイズを取得
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
          chunks: []
        });
      }

      const docEntry = documentMap.get(docId)!;

      // チャンクを追加（重複チェック）
      const existingChunk = docEntry.chunks.find(
        (c) => c.chunk_index === ref.chunk_index
      );

      if (!existingChunk) {
        docEntry.chunks.push({
          chunk_index: ref.chunk_index,
          content_preview: ref.content_preview || '',
          relevance_score: ref.score
        });
      }

      // このメッセージが参照するソースのインデックスを記録
      if (docEntry.sourceIndex === undefined) {
        docEntry.sourceIndex = allSources.length;
        allSources.push({
          ...docEntry.document as SourceDocument,
          chunks_used: docEntry.chunks
        });
      }

      if (!messageSourceIndices.includes(docEntry.sourceIndex)) {
        messageSourceIndices.push(docEntry.sourceIndex);
      }
    });

    // このメッセージのソースインデックスを記録
    if (messageSourceIndices.length > 0) {
      messageSourcesMap.set(message.id, messageSourceIndices);
    }
  });

  // 各ドキュメントのチャンクをスコア順にソート
  allSources.forEach((source) => {
    source.chunks_used.sort((a, b) => {
      const scoreA = a.relevance_score || 0;
      const scoreB = b.relevance_score || 0;
      return scoreB - scoreA; // 降順
    });
  });

  return {
    allSources,
    messageSourcesMap
  };
}

/**
 * テキストを指定文字数で切り詰める
 */
export function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) {
    return text;
  }
  return text.substring(0, maxLength) + '...';
}