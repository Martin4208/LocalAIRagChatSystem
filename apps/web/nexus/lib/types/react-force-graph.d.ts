declare module 'react-force-graph-3d' {
  import { FC, RefObject } from 'react';

  export interface ForceGraphMethods {
    d3Force: (name: string) => unknown;
    d3ReheatSimulation: () => void;
    pauseAnimation: () => void;
    resumeAnimation: () => void;
    cameraPosition: (position?: { x: number; y: number; z: number }) => void;
    zoomToFit: (duration?: number, padding?: number) => void;
  }

  export interface ForceGraphProps<NodeType = unknown, LinkType = unknown> {
    ref?: RefObject<ForceGraphMethods>;
    graphData: {
      nodes: NodeType[];
      links: LinkType[];
    };
    nodeLabel?: string | ((node: NodeType) => string);
    nodeAutoColorBy?: string;
    nodeVal?: string | number | ((node: NodeType) => number);
    linkLabel?: string | ((link: LinkType) => string);
    linkDirectionalArrowLength?: number;
    linkDirectionalArrowRelPos?: number;
    linkCurvature?: number;
    onNodeClick?: (node: NodeType) => void;
    onLinkClick?: (link: LinkType) => void;
    enableNodeDrag?: boolean;
    enableNavigationControls?: boolean;
  }

  const ForceGraph3D: FC<ForceGraphProps>;
  export default ForceGraph3D;
}