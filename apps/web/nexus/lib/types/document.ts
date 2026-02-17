export type DocumentStatus = 'uploaded' | 'processing' | 'processed' | 'failed';

export interface Document {
  id: string;
  workspace_id: string;
  directory_id?: string;
  file_id: string;
  name: string;
  tags: string[];
  status: DocumentStatus;
  processed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface DocumentChunk {
  id: string;
  document_id: string;
  chunk_index: number;
  content: string;
  qdrant_point_id?: string;
  created_at: string;
}

export interface ProcessDocumentOptions {
  chunk_size?: number;
  chunk_overlap?: number;
  force_reprocess?: boolean;
}
