import { apiClient } from './client';
import type {
  Graph,
  GraphListItem,
  GraphNode,
  GraphEdge,
  CreateGraphRequest,
  UpdateGraphRequest,
} from '@/lib/types/graph';

export const graphsApi = {
  /**
   * List all graphs in a workspace
   */
  async list(workspaceId: string): Promise<{ graphs: GraphListItem[] }> {
    return apiClient<{ graphs: GraphListItem[] }>(
      `/workspaces/${workspaceId}/graphs`,
      { method: 'GET' }
    );
  },

  /**
   * Create a new graph
   */
  async create(
    workspaceId: string,
    data: CreateGraphRequest
  ): Promise<{ graph: Graph }> {
    return apiClient<{ graph: Graph }>(
      `/workspaces/${workspaceId}/graphs`,
      {
        method: 'POST',
        body: JSON.stringify(data),
      }
    );
  },

  /**
   * Get complete graph with nodes and edges
   */
  async get(workspaceId: string, graphId: string): Promise<Graph> {
    return apiClient<Graph>(
      `/workspaces/${workspaceId}/graphs/${graphId}`,
      { method: 'GET' }
    );
  },

  /**
   * Update graph metadata
   */
  async update(
    workspaceId: string,
    graphId: string,
    data: UpdateGraphRequest
  ): Promise<{ graph: Graph }> {
    return apiClient<{ graph: Graph }>(
      `/workspaces/${workspaceId}/graphs/${graphId}`,
      {
        method: 'PUT',
        body: JSON.stringify(data),
      }
    );
  },

  /**
   * Delete a graph
   */
  async delete(workspaceId: string, graphId: string): Promise<void> {
    return apiClient<void>(
      `/workspaces/${workspaceId}/graphs/${graphId}`,
      { method: 'DELETE' }
    );
  },

  /**
   * Get nodes only (lightweight)
   */
  async getNodes(
    workspaceId: string,
    graphId: string
  ): Promise<{ nodes: GraphNode[] }> {
    return apiClient<{ nodes: GraphNode[] }>(
      `/workspaces/${workspaceId}/graphs/${graphId}/nodes`,
      { method: 'GET' }
    );
  },

  /**
   * Get edges only (lightweight)
   */
  async getEdges(
    workspaceId: string,
    graphId: string
  ): Promise<{ edges: GraphEdge[] }> {
    return apiClient<{ edges: GraphEdge[] }>(
      `/workspaces/${workspaceId}/graphs/${graphId}/edges`,
      { method: 'GET' }
    );
  },
};