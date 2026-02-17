'use client';

import { NODE_TYPE_STYLES } from '@/lib/constants/graph';

interface GraphLegendProps {
  nodeTypes: string[];
}

export function GraphLegend({ nodeTypes }: GraphLegendProps) {
  if (nodeTypes.length === 0) return null;

  return (
    <div className="absolute bottom-4 left-4 z-10 bg-white rounded-lg shadow-lg border p-3">
      <h3 className="text-xs font-medium mb-2 text-muted-foreground">Node Types</h3>
      <div className="space-y-1.5">
        {nodeTypes.map((type) => {
          const style = NODE_TYPE_STYLES[type] || NODE_TYPE_STYLES.default;
          return (
            <div key={type} className="flex items-center gap-2">
              <div
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: style.color }}
              />
              <span className="text-xs">{style.label || type}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}