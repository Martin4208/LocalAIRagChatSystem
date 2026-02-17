'use client';

import { useParams, useRouter } from 'next/navigation';
import { useWorkspaceContext } from '@/lib/providers/workspace-context';
import { useGraph, useUpdateGraph, useDeleteGraph } from '@/lib/hooks/use-graph';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ArrowLeft, Trash2, Edit2, Save, X } from 'lucide-react';
import { useState, useMemo } from 'react';
import type { ForceGraphNode, GraphFilter } from '@/lib/types/graph';
import { GraphVisualization } from '@/components/graph/graph-visualization';
import { GraphControls } from '@/components/graph/graph-controls';
import { NodeDetailSidebar } from '@/components/graph/node-detail-sidebar';
import { GraphLegend } from '@/components/graph/graph-legend';
import {
  convertToForceGraphData,
  filterGraphData,
  getUniqueNodeTypes,
  getUniqueEdgeTypes,
} from '@/lib/utils/graph-adapter';

export default function GraphDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { workspaceId } = useWorkspaceContext();
  const graphId = params.graphId as string;

  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [editedTitle, setEditedTitle] = useState('');
  const [viewMode, setViewMode] = useState<'2d' | '3d'>('2d');
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [hoveredNodeId, setHoveredNodeId] = useState<string | null>(null);
  const [filter, setFilter] = useState<GraphFilter>({
    nodeTypes: [],
    edgeTypes: [],
    searchQuery: '',
  });

  // Get graph data
  const { data: graph, isLoading, error } = useGraph(workspaceId, graphId);
  const updateGraphMutation = useUpdateGraph();
  const deleteGraphMutation = useDeleteGraph();

  // Convert and filter graph data
  const graphData = useMemo(() => {
    if (!graph) return { nodes: [], links: [] };
    const converted = convertToForceGraphData(graph);
    return filterGraphData(converted, filter);
  }, [graph, filter]);

  // Get available types for filters
  const availableNodeTypes = useMemo(
    () => (graph ? getUniqueNodeTypes(graph.nodes) : []),
    [graph]
  );

  const availableEdgeTypes = useMemo(
    () => (graph ? getUniqueEdgeTypes(graph.edges) : []),
    [graph]
  );

  // Get selected node and its connections
  const selectedNode = useMemo(
    () => graphData.nodes.find((n) => n.id === selectedNodeId),
    [graphData.nodes, selectedNodeId]
  );

  const connectedEdges = useMemo(() => {
    if (!selectedNodeId) return [];
    return graphData.links.filter((link) => {
      const sourceId = typeof link.source === 'string' ? link.source : link.source.id;
      const targetId = typeof link.target === 'string' ? link.target : link.target.id;
      return sourceId === selectedNodeId || targetId === selectedNodeId;
    });
  }, [graphData.links, selectedNodeId]);

  const handleBack = () => {
    router.push(`/workspaces/${workspaceId}/graph`);
  };

  const handleStartEdit = () => {
    setEditedTitle(graph?.title || '');
    setIsEditingTitle(true);
  };

  const handleSaveTitle = async () => {
    if (!editedTitle.trim()) {
      return;
    }

    try {
      await updateGraphMutation.mutateAsync({
        workspaceId,
        graphId,
        data: { title: editedTitle },
      });
      setIsEditingTitle(false);
    } catch (error) {
      console.error('Failed to update graph title:', error);
    }
  };

  const handleCancelEdit = () => {
    setIsEditingTitle(false);
    setEditedTitle('');
  };

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this graph?')) {
      return;
    }

    try {
      await deleteGraphMutation.mutateAsync({
        workspaceId,
        graphId,
      });
      router.push(`/workspaces/${workspaceId}/graph`);
    } catch (error) {
      console.error('Failed to delete graph:', error);
    }
  };

  const handleNodeClick = (node: ForceGraphNode) => {
    setSelectedNodeId(node.id);
  };

  const handleNodeHover = (node: ForceGraphNode | null) => {
    setHoveredNodeId(node?.id || null);
  };

  const handleBackgroundClick = () => {
    setSelectedNodeId(null);
  };

  const handleCloseSidebar = () => {
    setSelectedNodeId(null);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-muted-foreground">Loading graph...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4">
        <div className="text-destructive">
          Error loading graph: {String(error)}
        </div>
        <Button onClick={handleBack} className="mt-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Graphs
        </Button>
      </div>
    );
  }

  if (!graph) {
    return (
      <div className="p-4">
        <div className="text-muted-foreground">Graph not found</div>
        <Button onClick={handleBack} className="mt-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Graphs
        </Button>
      </div>
    );
  }

  return (
    <div className="h-screen w-full flex flex-col">
      {/* Header */}
      <div className="border-b p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4 flex-1">
            <Button variant="ghost" size="sm" onClick={handleBack}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            {isEditingTitle ? (
              <div className="flex items-center gap-2 flex-1 max-w-md">
                <Input
                  value={editedTitle}
                  onChange={(e) => setEditedTitle(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') handleSaveTitle();
                    if (e.key === 'Escape') handleCancelEdit();
                  }}
                  className="flex-1"
                  autoFocus
                />
                <Button
                  size="sm"
                  onClick={handleSaveTitle}
                  disabled={updateGraphMutation.isPending}
                >
                  <Save className="h-4 w-4" />
                </Button>
                <Button size="sm" variant="ghost" onClick={handleCancelEdit}>
                  <X className="h-4 w-4" />
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <h1 className="text-2xl font-bold">{graph.title}</h1>
                <Button variant="ghost" size="sm" onClick={handleStartEdit}>
                  <Edit2 className="h-4 w-4" />
                </Button>
              </div>
            )}
          </div>
          <div className="flex items-center gap-2">
            <div className="text-sm text-muted-foreground">
              {graphData.nodes.length} / {graph.nodes?.length || 0} nodes •{' '}
              {graphData.links.length} / {graph.edges?.length || 0} edges
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={handleDelete}
              disabled={deleteGraphMutation.isPending}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      {/* Graph Visualization */}
      <div className="flex-1 relative">
        {graph.nodes && graph.nodes.length > 0 ? (
          <>
            <GraphVisualization
              data={graphData}
              mode={viewMode}
              selectedNodeId={selectedNodeId}
              onNodeClick={handleNodeClick}
              onNodeHover={handleNodeHover}
              onBackgroundClick={handleBackgroundClick}
            />
            <GraphControls
              mode={viewMode}
              onModeChange={setViewMode}
              filter={filter}
              onFilterChange={setFilter}
              availableNodeTypes={availableNodeTypes}
              availableEdgeTypes={availableEdgeTypes}
            />
            <GraphLegend nodeTypes={availableNodeTypes} />
            {selectedNode && (
              <NodeDetailSidebar
                node={selectedNode}
                connectedEdges={connectedEdges}
                onClose={handleCloseSidebar}
              />
            )}
          </>
        ) : (
          <div className="h-full w-full flex items-center justify-center">
            <div className="text-center">
              <p className="text-muted-foreground text-lg">No nodes in this graph yet</p>
              <p className="text-sm text-muted-foreground mt-2">
                Add documents or run analysis to populate the graph
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}