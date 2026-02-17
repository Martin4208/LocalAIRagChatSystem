export interface Workspace {
    id: string,
    name: string,
    description?: string,
    settings?: Record<string, any>,
    created_at: string,
    updated_at: string,
    deleted_at?: string
}

export interface WorkspaceResponse {
    workspaces: Workspace[];
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/vi';

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    const response = await fetch(url, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options?.headers,
        },
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP${response.status}`);
    }
    
    return response.json();
}

export const workspaceApi = {
    list: async (): Promise<WorkspaceResponse> => {
        return fetchApi<WorkspaceResponse>('/workspaces');
    }
};
