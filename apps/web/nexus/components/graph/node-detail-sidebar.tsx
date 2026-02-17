'use client';

import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { X, ExternalLink, FileText } from 'lucide-react';
import type { ForceGraphNode, ForceGraphEdge } from '@/lib/types/graph';
import { format } from 'date-fns';

interface NodeDetailSidebarProps {
  node: ForceGraphNode;
  connectedEdges: ForceGraphEdge[];
  onClose: () => void;
}

export function NodeDetailSidebar({
  node,
  connectedEdges,
  onClose,
}: NodeDetailSidebarProps) {
  const incomingEdges = connectedEdges.filter((e) => {
    const targetId = typeof e.target === 'string' ? e.target : e.target.id;
    return targetId === node.id;
  });

  const outgoingEdges = connectedEdges.filter((e) => {
    const sourceId = typeof e.source === 'string' ? e.source : e.source.id;
    return sourceId === node.id;
  });

  return (
    <div className="absolute top-0 right-0 h-full w-96 bg-white border-l shadow-lg z-20 flex flex-col">
      {/* Header */}
      <div className="p-4 border-b flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <Badge variant="secondary">{node.node_type}</Badge>
            {node.source_type && (
              <Badge variant="outline" className="text-xs">
                {node.source_type}
              </Badge>
            )}
          </div>
          <h2 className="text-lg font-semibold truncate" title={node.label}>
            {node.label}
          </h2>
          <p className="text-xs text-muted-foreground mt-1">
            ID: {node.id.slice(0, 8)}...
          </p>
        </div>
        <Button variant="ghost" size="sm" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <ScrollArea className="flex-1">
        <div className="p-4 space-y-4">
          {/* Basic Info */}
          <div>
            <h3 className="text-sm font-medium mb-2">Information</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Created:</span>
                <span>{format(new Date(node.created_at), 'PPp')}</span>
              </div>
              {node.position && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Position:</span>
                  <span className="font-mono text-xs">
                    ({node.position.x?.toFixed(1)}, {node.position.y?.toFixed(1)}
                    {node.position.z && `, ${node.position.z.toFixed(1)}`})
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Source Link */}
          {node.source_id && (
            <>
              <Separator />
              <div>
                <h3 className="text-sm font-medium mb-2">Source</h3>
                <Button variant="outline" size="sm" className="w-full justify-start">
                  <FileText className="h-4 w-4 mr-2" />
                  View Source Document
                  <ExternalLink className="h-3 w-3 ml-auto" />
                </Button>
              </div>
            </>
          )}

          {/* Metadata */}
          {node.metadata && Object.keys(node.metadata).length > 0 && (
            <>
              <Separator />
              <div>
                <h3 className="text-sm font-medium mb-2">Metadata</h3>
                <div className="bg-muted/50 rounded p-3 text-xs font-mono">
                  <pre className="whitespace-pre-wrap break-words">
                    {JSON.stringify(node.metadata, null, 2)}
                  </pre>
                </div>
              </div>
            </>
          )}

          {/* Connections */}
          <Separator />
          <div>
            <h3 className="text-sm font-medium mb-2">Connections</h3>
            <div className="space-y-3">
              {/* Outgoing */}
              {outgoingEdges.length > 0 && (
                <div>
                  <p className="text-xs text-muted-foreground mb-2">
                    Outgoing ({outgoingEdges.length})
                  </p>
                  <div className="space-y-1">
                    {outgoingEdges.slice(0, 5).map((edge) => (
                      <div
                        key={edge.id}
                        className="flex items-center gap-2 text-sm p-2 rounded bg-muted/30"
                      >
                        <Badge variant="outline" className="text-xs">
                          {edge.edge_type}
                        </Badge>
                        <span className="text-xs truncate flex-1">
                          →{' '}
                          {typeof edge.target === 'string'
                            ? edge.target.slice(0, 8)
                            : edge.target.label}
                        </span>
                        {edge.weight && (
                          <span className="text-xs text-muted-foreground">
                            {edge.weight.toFixed(2)}
                          </span>
                        )}
                      </div>
                    ))}
                    {outgoingEdges.length > 5 && (
                      <p className="text-xs text-muted-foreground text-center py-1">
                        +{outgoingEdges.length - 5} more
                      </p>
                    )}
                  </div>
                </div>
              )}

              {/* Incoming */}
              {incomingEdges.length > 0 && (
                <div>
                  <p className="text-xs text-muted-foreground mb-2">
                    Incoming ({incomingEdges.length})
                  </p>
                  <div className="space-y-1">
                    {incomingEdges.slice(0, 5).map((edge) => (
                      <div
                        key={edge.id}
                        className="flex items-center gap-2 text-sm p-2 rounded bg-muted/30"
                      >
                        <Badge variant="outline" className="text-xs">
                          {edge.edge_type}
                        </Badge>
                        <span className="text-xs truncate flex-1">
                          ←{' '}
                          {typeof edge.source === 'string'
                            ? edge.source.slice(0, 8)
                            : edge.source.label}
                        </span>
                        {edge.weight && (
                          <span className="text-xs text-muted-foreground">
                            {edge.weight.toFixed(2)}
                          </span>
                        )}
                      </div>
                    ))}
                    {incomingEdges.length > 5 && (
                      <p className="text-xs text-muted-foreground text-center py-1">
                        +{incomingEdges.length - 5} more
                      </p>
                    )}
                  </div>
                </div>
              )}

              {outgoingEdges.length === 0 && incomingEdges.length === 0 && (
                <p className="text-sm text-muted-foreground text-center py-4">
                  No connections
                </p>
              )}
            </div>
          </div>
        </div>
      </ScrollArea>
    </div>
  );
}