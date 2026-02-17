import { apiClient } from './client';

export interface ChunkReference {
  chunk_index: number;
  score: number;
  content_preview: string;
  relevance_score?: number;
}

export interface SourceDocument {
  document_id: string;
  document_name: string;
  mime_type: string;
  size_bytes: number;
  chunks_used: ChunkReference[];
  created_at: string;
  updated_at: string;
}

export interface SourcesResponse {
  sources: SourceDocument[];
  total_documents: number;
  total_chunks: number;
}

export const sourcesApi = {
  getChatSources: async (
    workspaceId: string,
    chatId: string
  ): Promise<SourcesResponse> => {
    const response = await apiClient.get(
      `/workspaces/${workspaceId}/chats/${chatId}/sources`
    );
    return response.data;
  },

  getAnalysisSources: async (
    workspaceId: string,
    analysisId: string
  ): Promise<SourcesResponse> => {
    const response = await apiClient.get(
      `/workspaces/${workspaceId}/analyses/${analysisId}/sources`
    );
    return response.data;
  },
};