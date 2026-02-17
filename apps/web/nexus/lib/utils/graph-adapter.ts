// ========================================
// Graph Data Adapter
// ========================================

import type {
  Graph,
  GraphNode,
  GraphEdge,
  ForceGraphNode,
  ForceGraphEdge,
  ForceGraphData,
  GraphFilter,
} from '@/lib/types/graph';
import { NODE_TYPE_STYLES, EDGE_TYPE_STYLES, DEFAULT_NODE_SIZE } from '@/lib/constants/graph';

/**
 * Convert API Graph to ForceGraph data structure
 */
export function convertToForceGraphData(graph: Graph): ForceGraphData {
  const nodes: ForceGraphNode[] = graph.nodes.map((node) => ({
    ...node,
    // Initialize position if not set
    x: node.position?.x,
    y: node.position?.y,
    z: node.position?.z ?? undefined,
    // Apply styling
    color: getNodeColor(node),
    size: getNodeSize(node),
  }));

  const links: ForceGraphEdge[] = graph.edges.map((edge) => ({
    ...edge,
    source: edge.from_node_id,
    target: edge.to_node_id,
    color: getEdgeColor(edge),
    width: getEdgeWidth(edge),
  }));

  return { nodes, links };
}

/**
 * Get node color based on type and custom style
 */
function getNodeColor(node: GraphNode): string {
  // Custom style takes precedence
  if (node.style && typeof node.style.color === 'string') {
    return node.style.color;
  }

  // Fallback to type-based color
  const typeStyle = NODE_TYPE_STYLES[node.node_type];
  return typeStyle?.color || NODE_TYPE_STYLES.default.color;
}

/**
 * Get node size based on metadata or default
 */
function getNodeSize(node: GraphNode): number {
  if (node.style && typeof node.style.size === 'number') {
    return node.style.size;
  }

  const typeStyle = NODE_TYPE_STYLES[node.node_type];
  return typeStyle?.size || DEFAULT_NODE_SIZE;
}

/**
 * Get edge color based on type, weight, and confidence
 */
function getEdgeColor(edge: GraphEdge): string {
  // Custom style
  if (edge.style && typeof edge.style.color === 'string') {
    return edge.style.color;
  }

  // Low confidence = faded
  if (edge.confidence !== null && edge.confidence !== undefined && edge.confidence < 0.5) {
    return '#cbd5e1'; // slate-300
  }

  const typeStyle = EDGE_TYPE_STYLES[edge.edge_type];
  return typeStyle?.color || EDGE_TYPE_STYLES.default.color;
}

/**
 * Get edge width based on weight
 */
function getEdgeWidth(edge: GraphEdge): number {
  if (edge.style && typeof edge.style.width === 'number') {
    return edge.style.width;
  }

  // Scale by weight (0.5 - 3.0)
  const weight = edge.weight ?? 0.5;
  return 0.5 + weight * 2.5;
}

/**
 * Filter graph data based on criteria
 */
export function filterGraphData(
  data: ForceGraphData,
  filter: GraphFilter
): ForceGraphData {
  let filteredNodes = [...data.nodes];
  let filteredLinks = [...data.links];

  // Filter by node types
  if (filter.nodeTypes && filter.nodeTypes.length > 0) {
    filteredNodes = filteredNodes.filter((node) =>
      filter.nodeTypes!.includes(node.node_type)
    );
  }

  // Filter by search query (label contains)
  if (filter.searchQuery && filter.searchQuery.trim()) {
    const query = filter.searchQuery.toLowerCase();
    filteredNodes = filteredNodes.filter((node) =>
      node.label.toLowerCase().includes(query)
    );
  }

  // Get valid node IDs
  const nodeIds = new Set(filteredNodes.map((n) => n.id));

  // Filter edges: both nodes must exist
  filteredLinks = filteredLinks.filter((link) => {
    const sourceId = typeof link.source === 'string' ? link.source : link.source.id;
    const targetId = typeof link.target === 'string' ? link.target : link.target.id;
    return nodeIds.has(sourceId) && nodeIds.has(targetId);
  });

  // Filter by edge types
  if (filter.edgeTypes && filter.edgeTypes.length > 0) {
    filteredLinks = filteredLinks.filter((link) =>
      filter.edgeTypes!.includes(link.edge_type)
    );
  }

  // Filter by weight
  if (filter.minWeight !== undefined) {
    filteredLinks = filteredLinks.filter(
      (link) => (link.weight ?? 0) >= filter.minWeight!
    );
  }

  // Filter by confidence
  if (filter.minConfidence !== undefined) {
    filteredLinks = filteredLinks.filter(
      (link) => (link.confidence ?? 0) >= filter.minConfidence!
    );
  }

  return { nodes: filteredNodes, links: filteredLinks };
}

/**
 * Get unique node types from graph
 */
export function getUniqueNodeTypes(nodes: GraphNode[]): string[] {
  return Array.from(new Set(nodes.map((n) => n.node_type))).sort();
}

/**
 * Get unique edge types from graph
 */
export function getUniqueEdgeTypes(edges: GraphEdge[]): string[] {
  return Array.from(new Set(edges.map((e) => e.edge_type))).sort();
}