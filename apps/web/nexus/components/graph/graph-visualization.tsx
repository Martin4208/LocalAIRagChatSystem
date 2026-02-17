'use client';

import { useRef, useCallback, useEffect } from 'react';
import dynamic from 'next/dynamic';
import type { ForceGraphMethods } from 'react-force-graph-2d';
import type { ForceGraphData, ForceGraphNode, ForceGraphEdge } from '@/lib/types/graph';
import { FORCE_GRAPH_CONFIG } from '@/lib/constants/graph';

// Dynamic import to avoid SSR issues
const ForceGraph2D = dynamic(() => import('react-force-graph-2d'), { ssr: false });
const ForceGraph3D = dynamic(() => import('react-force-graph-3d'), { ssr: false });

interface GraphVisualizationProps {
  data: ForceGraphData;
  mode: '2d' | '3d';
  selectedNodeId: string | null;
  onNodeClick: (node: ForceGraphNode) => void;
  onNodeHover: (node: ForceGraphNode | null) => void;
  onBackgroundClick?: () => void;
}

export function GraphVisualization({
  data,
  mode,
  selectedNodeId,
  onNodeClick,
  onNodeHover,
  onBackgroundClick,
}: GraphVisualizationProps) {
  const graphRef = useRef<ForceGraphMethods>();

  // Zoom to fit on data change
  useEffect(() => {
    if (graphRef.current && data.nodes.length > 0) {
      setTimeout(() => {
        graphRef.current?.zoomToFit(400, 50);
      }, 100);
    }
  }, [data.nodes.length]);

  // Node canvas rendering
  const paintNode = useCallback(
    (node: ForceGraphNode, ctx: CanvasRenderingContext2D, globalScale: number) => {
      const label = node.label;
      const fontSize = 12 / globalScale;
      const nodeSize = node.size || 5;
      const isSelected = node.id === selectedNodeId;

      // Draw node circle
      ctx.beginPath();
      ctx.arc(node.x!, node.y!, nodeSize, 0, 2 * Math.PI);
      ctx.fillStyle = node.color || '#64748b';
      ctx.fill();

      // Selection ring
      if (isSelected) {
        ctx.strokeStyle = '#3b82f6';
        ctx.lineWidth = 2 / globalScale;
        ctx.stroke();
      }

      // Draw label
      ctx.font = `${fontSize}px Sans-Serif`;
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillStyle = '#1e293b';
      ctx.fillText(label, node.x!, node.y! + nodeSize + fontSize);
    },
    [selectedNodeId]
  );

  // Link canvas rendering
  const paintLink = useCallback(
    (link: ForceGraphEdge, ctx: CanvasRenderingContext2D, globalScale: number) => {
      const source = typeof link.source === 'object' ? link.source : null;
      const target = typeof link.target === 'object' ? link.target : null;

      if (!source || !target) return;

      const lineWidth = (link.width || 1) / globalScale;
      ctx.strokeStyle = link.color || '#94a3b8';
      ctx.lineWidth = lineWidth;

      // Draw line
      ctx.beginPath();
      ctx.moveTo(source.x!, source.y!);
      ctx.lineTo(target.x!, target.y!);
      ctx.stroke();

      // Draw arrow if directed
      if (link.is_directed) {
        const arrowLength = 8 / globalScale;
        const angle = Math.atan2(target.y! - source.y!, target.x! - source.x!);
        const targetNodeSize = target.size || 5;

        // Arrow tip position (on edge of target node)
        const arrowX = target.x! - targetNodeSize * Math.cos(angle);
        const arrowY = target.y! - targetNodeSize * Math.sin(angle);

        ctx.save();
        ctx.translate(arrowX, arrowY);
        ctx.rotate(angle);
        ctx.beginPath();
        ctx.moveTo(0, 0);
        ctx.lineTo(-arrowLength, arrowLength / 2);
        ctx.lineTo(-arrowLength, -arrowLength / 2);
        ctx.closePath();
        ctx.fillStyle = link.color || '#94a3b8';
        ctx.fill();
        ctx.restore();
      }
    },
    []
  );

  const commonProps = {
    graphData: data,
    nodeId: 'id',
    nodeLabel: 'label',
    nodeVal: (node: ForceGraphNode) => node.size || 5,
    nodeColor: (node: ForceGraphNode) => node.color || '#64748b',
    linkSource: 'source',
    linkTarget: 'target',
    linkWidth: (link: ForceGraphEdge) => link.width || 1,
    linkColor: (link: ForceGraphEdge) => link.color || '#94a3b8',
    linkDirectionalArrowLength: 6,
    linkDirectionalArrowRelPos: 1,
    onNodeClick: (node: ForceGraphNode) => onNodeClick(node),
    onNodeHover: (node: ForceGraphNode | null) => onNodeHover(node),
    onBackgroundClick: onBackgroundClick,
    ...FORCE_GRAPH_CONFIG[mode],
  };

  if (mode === '2d') {
    return (
      <ForceGraph2D
        ref={graphRef as any}
        {...commonProps}
        nodeCanvasObject={paintNode}
        linkCanvasObject={paintLink}
        backgroundColor="#ffffff"
      />
    );
  }

  return (
    <ForceGraph3D
      ref={graphRef as any}
      {...commonProps}
      backgroundColor="#ffffff"
      showNavInfo={false}
    />
  );
}