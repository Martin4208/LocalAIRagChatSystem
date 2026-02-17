import { apiClient } from './client';
import type { Workspace } from '@/types/domain';

// 単一Workspace取得
export async function getWorkspace(id: string): Promise<Workspace> {
  return apiClient<Workspace>(`/workspaces/${id}`);
}

// 全Workspace取得
export async function getWorkspaces(): Promise<{ workspaces: Workspace[] }> {
  return apiClient<{ workspaces: Workspace[] }>('/workspaces');
}

// Workspace作成
export async function createWorkspace(data: {
  name: string;
  description?: string;
}): Promise<Workspace> {
  return apiClient<Workspace>('/workspaces', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Workspace更新
export async function updateWorkspace(
  id: string,
  data: {
    name?: string;
    description?: string;
  }
): Promise<Workspace> {
  return apiClient<Workspace>(`/workspaces/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

// Workspace削除
export async function deleteWorkspace(id: string): Promise<void> {
  return apiClient<void>(`/workspaces/${id}`, {
    method: 'DELETE',
  });
}