// ========================================
// Graph Styling Constants
// ========================================

export const DEFAULT_NODE_SIZE = 5;
export const DEFAULT_EDGE_WIDTH = 1;

/**
 * Node type color mappings
 * Using Tailwind color palette for consistency
 */
export const NODE_TYPE_STYLES: Record<
  string,
  { color: string; size: number; label?: string }
> = {
  // Common entity types
  person: {
    color: '#3b82f6', // blue-500
    size: 6,
    label: 'Person',
  },
  organization: {
    color: '#8b5cf6', // violet-500
    size: 7,
    label: 'Organization',
  },
  location: {
    color: '#10b981', // green-500
    size: 5,
    label: 'Location',
  },
  document: {
    color: '#f59e0b', // amber-500
    size: 8,
    label: 'Document',
  },
  concept: {
    color: '#ec4899', // pink-500
    size: 5,
    label: 'Concept',
  },
  event: {
    color: '#ef4444', // red-500
    size: 6,
    label: 'Event',
  },
  product: {
    color: '#06b6d4', // cyan-500
    size: 5,
    label: 'Product',
  },
  
  // Document-derived types
  document_chunk: {
    color: '#fbbf24', // amber-400
    size: 4,
    label: 'Chunk',
  },
  keyword: {
    color: '#a855f7', // purple-500
    size: 3,
    label: 'Keyword',
  },
  
  // Default fallback
  default: {
    color: '#64748b', // slate-500
    size: DEFAULT_NODE_SIZE,
    label: 'Unknown',
  },
};

/**
 * Edge type color mappings
 */
export const EDGE_TYPE_STYLES: Record<
  string,
  { color: string; width?: number; label?: string }
> = {
  // Relationships
  works_for: {
    color: '#3b82f6', // blue-500
    label: 'Works For',
  },
  located_in: {
    color: '#10b981', // green-500
    label: 'Located In',
  },
  mentions: {
    color: '#f59e0b', // amber-500
    width: 1,
    label: 'Mentions',
  },
  related_to: {
    color: '#8b5cf6', // violet-500
    label: 'Related To',
  },
  contains: {
    color: '#06b6d4', // cyan-500
    label: 'Contains',
  },
  references: {
    color: '#ec4899', // pink-500
    label: 'References',
  },
  
  // Semantic relationships
  similar_to: {
    color: '#a855f7', // purple-500
    width: 0.8,
    label: 'Similar To',
  },
  contradicts: {
    color: '#ef4444', // red-500
    width: 1.5,
    label: 'Contradicts',
  },
  supports: {
    color: '#22c55e', // green-500
    width: 1.2,
    label: 'Supports',
  },
  
  // Default
  default: {
    color: '#94a3b8', // slate-400
    width: DEFAULT_EDGE_WIDTH,
    label: 'Relationship',
  },
};

/**
 * Force graph physics config
 */
export const FORCE_GRAPH_CONFIG = {
  '2d': {
    cooldownTicks: 100,
    cooldownTime: 15000,
    d3AlphaDecay: 0.0228,
    d3VelocityDecay: 0.4,
    warmupTicks: 0,
    // Link force
    linkDistance: 100,
    linkStrength: 1,
    // Charge force (repulsion)
    chargeStrength: -120,
  },
  '3d': {
    cooldownTicks: 100,
    cooldownTime: 15000,
    d3AlphaDecay: 0.0228,
    d3VelocityDecay: 0.4,
    warmupTicks: 0,
    linkDistance: 150,
    linkStrength: 1,
    chargeStrength: -180,
  },
};