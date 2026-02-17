import { apiClient } from './client'
import type { DocumentWithFile, DirectoryTree } from '@/types/domain'

export interface FileMetadata {
  id: string;
  fileName: string;
  mimeType: string;
  sizeBytes: number;
  sha256Hash?: string;
  createdAt: string;
  tags?: string[];
}

export interface FileUploadResponse {
  id: string;
  fileName: string;
  mimeType: string;
  sizeBytes: number;
  sha256Hash: string;
  createdAt: string;
}

export const files = {
  /**
   * ファイル一覧取得
   */
  list: async (
    workspaceId: string,
    params?: {
      directoryId?: string;
      limit?: number;
      offset?: number;
    }
  ): Promise<{ files: FileMetadata[]; total: number }> => {
    const queryParams = new URLSearchParams();
    if (params?.directoryId) queryParams.append('directoryId', params.directoryId);
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const queryString = queryParams.toString();
    const endpoint = queryString 
      ? `/workspaces/${workspaceId}/files?${queryString}`
      : `/workspaces/${workspaceId}/files`;

    return apiClient<{ files: FileMetadata[]; total: number }>(endpoint);
  },

  /**
   * ファイルアップロード
   */
  upload: async (
    workspaceId: string,
    file: File,
    directoryId?: string,
    tags?: string[]
  ): Promise<FileUploadResponse> => {
    const formData = new FormData();
    formData.append('file', file);

    if (directoryId) {
      formData.append('directoryId', directoryId);
    }

    if (tags && tags.length > 0) {
      formData.append('tags', JSON.stringify(tags));
    }

    return  apiClient<FileUploadResponse>(
      `/workspaces/${workspaceId}/files/upload`,
      {
        method: 'POST',
        body: formData,
      }
    );
  },

  /**
   * ファイル取得
   */
  get: async (
    workspaceId: string,
    fileId: string
  ): Promise<FileMetadata> => {
    return apiClient<FileMetadata>(
      `/workspaces/${workspaceId}/files/${fileId}`
    );
  },

  /**
   * ファイル削除
   */
  delete: async (workspaceId: string, fileId: string): Promise<void> => {
    return apiClient<void>(
      `/workspaces/${workspaceId}/files/${fileId}`,
      {
        method: 'DELETE',
      }
    );
  },

  /**
   * ファイルダウンロード
   */
  download: async (
    workspaceId: string,
    fileId: string
  ): Promise<Blob> => {
    const url = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/workspaces/${workspaceId}/files/${fileId}/download`;
    
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`Download failed: ${response.statusText}`);
    }
    
    return response.blob();
  },
};

export async function getWorkspaceDocuments(
    workspaceId: string
): Promise<DocumentWithFile[]> {
    return apiClient<DocumentWithFile[]>(
        `/workspaces/${workspaceId}/documents`
    )
}

export async function getDirectoryTree(
    workspaceId: string
): Promise<DirectoryTree[]> {
    return apiClient<DirectoryTree[]>(
        `/workspaces/${workspaceId}/directories/tree` 
    )
}
