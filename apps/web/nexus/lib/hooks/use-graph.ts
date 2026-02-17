import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { graphsApi } from '@/lib/api/graphs';
import type {
  Graph,
  GraphListItem,
  CreateGraphRequest,
  UpdateGraphRequest,
} from '@/lib/types/graph';

const MOCK_GRAPH = {
  id: '123',
  workspace_id: 'abc',
  title: 'Test Graph',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  nodes: [
    {
      id: 'node1',
      graph_id: '123',
      label: 'Alice',
      node_type: 'person',
      created_at: new Date().toISOString(),
    },
    {
      id: 'node2',
      graph_id: '123',
      label: 'Company X',
      node_type: 'organization',
      created_at: new Date().toISOString(),
    },
  ],
  edges: [
    {
      id: 'edge1',
      graph_id: '123',
      from_node_id: 'node1',
      to_node_id: 'node2',
      edge_type: 'works_for',
      is_directed: true,
      weight: 0.9,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
  ],
};


export function useGraphs(workspaceId: string | undefined) {
  return useQuery({
    queryKey: ['graphs', workspaceId],
    queryFn: () => graphsApi.list(workspaceId!),
    enabled: !!workspaceId,
    select: (data) => data.graphs,
  });
}

export function useGraph(workspaceId: string | undefined, graphId: string | undefined) {
  return useQuery({
    queryKey: ['graphs', workspaceId, graphId],
    queryFn: () => graphsApi.get(workspaceId!, graphId!),
    enabled: !!workspaceId && !!graphId,
  });
}

export function useCreateGraph() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ workspaceId, data }: { workspaceId: string; data: CreateGraphRequest }) =>
      graphsApi.create(workspaceId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['graphs', variables.workspaceId] });
    },
  });
}

export function useUpdateGraph() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workspaceId,
      graphId,
      data,
    }: {
      workspaceId: string;
      graphId: string;
      data: UpdateGraphRequest;
    }) => graphsApi.update(workspaceId, graphId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['graphs', variables.workspaceId, variables.graphId],
      });
      queryClient.invalidateQueries({ queryKey: ['graphs', variables.workspaceId] });
    },
  });
}

export function useDeleteGraph() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ workspaceId, graphId }: { workspaceId: string; graphId: string }) =>
      graphsApi.delete(workspaceId, graphId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['graphs', variables.workspaceId] });
    },
  });
}