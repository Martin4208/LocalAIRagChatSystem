// ========================================
// Graph Type Definitions
// ========================================

// API Response Types (from OpenAPI schema)
export interface GraphListItem {
  id: string;
  workspace_id: string;
  title: string;
  graph_type?: string | null;
  created_at: string;
  updated_at: string;
}

export interface GraphNode {
  id: string;
  graph_id: string;
  label: string;
  node_type: string;
  source_type?: string | null;
  source_id?: string | null;
  position?: {
    x: number;
    y: number;
    z?: number | null;
  } | null;
  style?: Record<string, unknown> | null;
  metadata?: Record<string, unknown> | null;
  created_at: string;
}

export interface GraphEdge {
  id: string;
  graph_id: string;
  from_node_id: string;
  to_node_id: string;
  edge_type: string;
  is_directed: boolean;
  weight?: number | null;
  confidence?: number | null;
  style?: Record<string, unknown> | null;
  metadata?: Record<string, unknown> | null;
  created_at: string;
  updated_at: string;
}

export interface Graph extends GraphListItem {
  layout_config?: Record<string, unknown> | null;
  nodes: GraphNode[];
  edges: GraphEdge[];
}

// UI-specific types for react-force-graph
export interface ForceGraphNode extends GraphNode {
  // Additional properties for force-graph
  x?: number;
  y?: number;
  z?: number;
  vx?: number;
  vy?: number;
  vz?: number;
  fx?: number | null;
  fy?: number | null;
  fz?: number | null;
  // Visual properties
  color?: string;
  size?: number;
}

export interface ForceGraphEdge extends Omit<GraphEdge, 'from_node_id' | 'to_node_id'> {
  source: string | ForceGraphNode;
  target: string | ForceGraphNode;
  color?: string;
  width?: number;
}

export interface ForceGraphData {
  nodes: ForceGraphNode[];
  links: ForceGraphEdge[];
}

// Filter & Search
export interface GraphFilter {
  nodeTypes?: string[];
  edgeTypes?: string[];
  minWeight?: number;
  minConfidence?: number;
  searchQuery?: string;
}

// UI State
export interface GraphViewState {
  mode: '2d' | '3d';
  selectedNodeId: string | null;
  filter: GraphFilter;
  hoveredNodeId: string | null;
}