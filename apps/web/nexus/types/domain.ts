export interface Workspace {
    id: string;
    name: string;
    description?: string;
    settings?: Record<string, unknown>;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface Chat {
    id: string;
    workspace_id: string;
    title: string;
    filter_config?: Record<string, unknown>;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

// --- Files ---
export interface File {
    id: string;
    sha256_hash: string;
    mime_type: string;
    size_bytes: number;  // ← bigint だけど、JSON では number
    original_filename?: string;
    minio_bucket: string;
    minio_key: string;
    preview_content?: string;
    created_at: string;
    last_accessed_at?: string;
}

export interface Document {
    d: string;
    workspace_id: string;
    directory_id?: string;
    file_id: string;
    name: string;
    tags?: string[];  // ← PostgreSQL の array
    metadata?: Record<string, unknown>;  // ← JSONB
    status: 'uploaded' | 'processing' | 'ready' | 'error';  // 'uploaded' | 'processing' | 'completed' など
    processed_at?: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface Directory {
    id: string;
    workspace_id: string;
    parent_id?: string;
    name: string;
    created_at: string;
    updated_at: string;
}

export interface GetFilesResponse {
    documents: DocumentWithFile[]; // Document + Fileの結合データ
    total: number;
    page?: number;
    limit?: number;
}

export interface DocumentWithFile extends Document {
    file: File;
}

export interface DirectoryTree extends Directory {
    children?: DirectoryTree[]; // 再帰的な構造
    document_count?: number;
}

export interface UploadFileRequest {
    file: File;
    directory_id?: string;
    tags?: string[];
}

export interface UploadProgress {
    file_id: string;
    filename: string;
    progress: number;
    status: 'pending' | 'uploading' | 'processing' | 'complete' | 'error';
    error_message?: string;
}