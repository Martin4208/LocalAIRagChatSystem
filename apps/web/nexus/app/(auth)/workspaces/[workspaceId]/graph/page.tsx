'use client';

import { useWorkspaceContext } from '@/lib/providers/workspace-context';
import { useGraphs, useCreateGraph, useDeleteGraph } from '@/lib/hooks/use-graph';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Trash2, Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';

export default function GraphPage() {
  const { workspaceId } = useWorkspaceContext();
  const router = useRouter();

  // Get list of graphs
  const { data: graphs, isLoading, error } = useGraphs(workspaceId);

  const createGraphMutation = useCreateGraph();
  const deleteGraphMutation = useDeleteGraph();

  // // For now, select the first graph (you can add graph selector UI later)
  // const selectedGraphId = graphs?.[0]?.id;

  // Get full graph data
  // const {
  //   data: graph,
  //   isLoading: isLoadingGraph,
  //   error,
  // } = useGraph(workspaceId, selectedGraphId);

  const handleCreateGraph = async () => {
    console.log('Creating Graph');
    try {
      await createGraphMutation.mutateAsync({
        workspaceId,
        data: {
          title: `New Graph ${new Date().toLocaleTimeString()}`
        }
      });
      console.log('✅ Graph created successfully');
      if (result.graph?.id) {
        router.push(`/workspaces/${workspaceId}/graph/${result.graph.id}`);
      }
    } catch (error) {
      console.log('Failed to create chat:', error)

      if (error?.response) {
        console.error('API Response:', error.response);
        console.error('Status:', error.response.status);
        console.error('Data:', error.response.data);
      }
    }
  }

  const handleDeleteGraph = async (graphId: string) => {
    if (!confirm('Are you sure you want to delete this graph?')) {
      return;
    }

    try {
      await deleteGraphMutation.mutateAsync({
        workspaceId,
        graphId
      });
      console.log('Successfully deleted graph');
    } catch (error) {
      console.log('Failed to delete graph:', error);
    }
  }

  const handleGraphClick = (graphId: string) => {
    router.push(`/workspaces/${workspaceId}/graph/${graphId}`);
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-muted-foreground">Loading graph...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-destructive">
        Error loading graph: {String(error)}
      </div>
    );
  }

  return (
    <div className="h-screen w-full">
      {/* Header */}
      <div className="border-b p-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold">Graphs</h1>
            <p className="text-sm text-muted-foreground mt-1">
              {graphs?.length || 0} graphs
            </p>
          </div>
          <Button
            onClick={handleCreateGraph}
            disabled={createGraphMutation.isPending}
            className="flex items-center gap-2"
          >
            <Plus className="h-4 w-4" />
            {createGraphMutation.isPending ? 'Creating...' : 'New Graph'}
          </Button>
        </div>
      </div>

      {/** グラフ一覧 */}
      <div>
        {graphs && graphs.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {graphs.map((graph) => (
              <Card
                key={graph.id}
                className="p-4 hover:shadow-lg transition-shadow cursor-pointer"
                onClick={() => handleGraphClick(graph.id)}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="font-semibold text-lg">{graph.title}</h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      {graph.graph_type || 'Knowledge Graph'}
                    </p>
                    <p>
                      Created: {new Date(graph.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeleteGraph(graph.id)
                    }}
                    disabled={deleteGraphMutation.isPending}
                    className="text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-4 w-4"/>
                  </Button>
                </div>
              </Card>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-[60vh] text-center">
            <p className="text-muted-foreground text-lg">No graphs yet</p>
            <p className="text-sm text-muted-foreground mt-2">
              Create your first knowledge graph to get started
            </p>
            <Button
              onClick={handleCreateGraph}
              disabled={createGraphMutation.isPending}
              className="mt-4"
            >
              <Plus className="h-4 w-4 mr-2" />
              Create Graph
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}