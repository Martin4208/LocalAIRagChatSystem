// apps/web/nexus/lib/api/documents.ts

import { apiClient } from './client';

export interface ProcessDocumentRequest {
  chunk_size?: number;
  chunk_overlap?: number;
  force_reprocess?: boolean;
}

export interface DocumentStatus {
  status: 'uploaded' | 'processing' | 'processed' | 'failed';
  progress?: {
    current_step?: string;
    percentage?: number;
    chunks_created?: number;
  } | null;
  error?: {
    message?: string;
    code?: string;
  } | null;
  processed_at?: string | null;
}

export interface DocumentChunk {
  id: string;
  chunk_index: number;
  content: string;
  created_at: string;
}

export interface GetChunksResponse {
  chunks: DocumentChunk[];
  total: number;
  page: number;
  limit: number;
}

export const documents = {
  /**
   * ドキュメントの処理を開始
   */
  process: async (
    workspaceId: string,
    documentId: string,
    options?: ProcessDocumentRequest
  ): Promise<{ status: string; message: string }> => {
    const response = await apiClient.post(
      `/workspaces/${workspaceId}/documents/${documentId}/process`,
      options || {}
    );
    return response.data;
  },

  /**
   * ドキュメントの処理状況を取得
   */
  getStatus: async (
    workspaceId: string,
    documentId: string
  ): Promise<DocumentStatus> => {
    const response = await apiClient.get(
      `/workspaces/${workspaceId}/documents/${documentId}/status`
    );
    return response.data;
  },

  /**
   * ドキュメントのチャンク一覧を取得
   */
  getChunks: async (
    workspaceId: string,
    documentId: string,
    page: number = 1,
    limit: number = 20
  ): Promise<GetChunksResponse> => {
    const response = await apiClient.get(
      `/workspaces/${workspaceId}/documents/${documentId}/chunks`,
      {
        params: { page, limit }
      }
    );
    return response.data;
  }
};