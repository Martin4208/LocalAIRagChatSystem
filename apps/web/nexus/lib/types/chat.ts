// lib/types/chat.ts

export type ChatRole = "user" | "assistant" | "system"

export interface FilterConfig {
    document_ids?: string[];
    directory_ids?: string[];
    tags?: string[];
}

export interface DocumentReference {
    document_id: string;
    document_name?: string;
    chunk_index: number;
    score: number;
    content_preview?: string;
}

export interface Chat {
    id: string;
    workspace_id: string;
    title: string;
    filter_config?: FilterConfig;
    message_count?: number;
    last_message_at?: string;
    created_at: string;
    updated_at: string;
}

export interface ChatMessage {
    id: string;
    chat_id: string;
    role: ChatRole;
    content: string;
    message_index: number;
    documentRefs?: DocumentReference[];
    created_at: string;
}

export interface CreateChatRequest {
    title: string;
    filter_config?: FilterConfig;
}

export interface ListChatResponse {
    chats: Chat[];
    total: number;
}

export interface GetChatResponse {
    chat: Chat;
    messages: ChatMessage[];
}

export interface SendMessageRequest {
    content: string;
}

export interface SendMessageResponse {
    user_message: ChatMessage;
    assistant_message: ChatMessage;
}

export interface SourceMapping {
  index: number;              // [1], [2], [3]...
  documentId: string;
  documentName: string;
  chunkIndex: number;
  score: number;
  contentPreview: string;
}

// メッセージ内のソース番号を抽出した結果
export interface MessageSources {
  messageId: string;
  sourceIndices: number[];    // このメッセージで参照されている番号
}
