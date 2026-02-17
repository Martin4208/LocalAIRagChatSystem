// apps/web/nexus/components/shared/sources-panel.tsx

'use client';

import { useSourcePanelStore } from '@/lib/stores/use-source-panel-store';
import { SourceDocument } from '@/lib/types/sources';
import { FileText, ChevronRight } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card } from '@/components/ui/card';
import { LAYOUT } from '@/lib/constants/layout';
import { useLayoutStore } from '@/lib/stores/use-side-pane-store';

interface SourcesPanelProps {
    workspaceId: string;
    chatId: string;
    sourcesList: SourceDocument[]; // ← Chat Pageから渡される
}

export function SourcesPanel({ workspaceId, chatId, sourcesList }: SourcesPanelProps) {
    const selectedSourceIndex = useSourcePanelStore((state) => state.selectedSourceIndex);
    const selectSource = useSourcePanelStore((state) => state.selectSource);
    const isSourcePanelOpen = useLayoutStore((state) => state.isSourcePanelOpen);

    const sourceWidth = isSourcePanelOpen 
        ? LAYOUT.source.width.open
        : LAYOUT.source.width.closed

    return (
        <aside
            className="h-full flex-shrink-0 border-l bg-background overflow-hidden
                        transition-[width] duration-300 ease-in-out"
            style={{ width: sourceWidth }}
        >
            {sourcesList.length === 0 ? (
                <div className="p-8 text-center">
                    <FileText className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                    <p className="text-sm text-gray-500">No sources referenced</p>
                </div>
            ): (
                <div 
                    className="space-y-4 p-4 overflow-y-auto h-full"
                >

                    <div className="flex items-center justify-between">
                        <h2 className="text-lg font-semibold">Sources</h2>
                        <span className="text-sm text-gray-500">
                            {sourcesList.length} {sourcesList.length === 1 ? 'document' : 'documents'}
                        </span>
                    </div>

                    <div className="space-y-3">
                        {sourcesList.map((source, index) => {
                            const isSelected = selectedSourceIndex === index;

                            return (
                                <Card
                                    key={source.document_id}
                                    className={cn(
                                        'p-4 cursor-pointer transition-all',
                                        isSelected && 'bg-blue-50 border-blue-400'
                                    )}
                                    onClick={() => selectSource(index)}
                                >
                                    <div className="flex items-start gap-3">
                                        <span className={cn(
                                            "flex-shrink-0 font-mono text-xs px-1.5 py-0.5 rounded",
                                            isSelected ? "bg-blue-600 text-white" : "bg-gray-200"
                                        )}>
                                            [{index}]
                                        </span>

                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-start gap-2">
                                                <FileText className="h-4 w-4 mt-0.5 text-gray-400" />
                                                <div>
                                                    <h3 className="font-medium text-sm truncate">
                                                        {source.document_name}
                                                    </h3>
                                                    <p className="text-xs text-gray-500">
                                                        {source.chunks_used.length} chunk{source.chunks_used.length !== 1 ? 's' : ''}
                                                    </p>
                                                </div>
                                            </div>

                                            {/* Chunks */}
                                            <div className="mt-3 space-y-2">
                                                {source.chunks_used.map((chunk) => (
                                                    <div 
                                                        key={chunk.chunk_index}
                                                        className="bg-gray-50 rounded p-2 text-xs"
                                                    >
                                                        <div className="flex items-center justify-between mb-1">
                                                            <span className="font-medium">Chunk {chunk.chunk_index}</span>
                                                            <span className="text-gray-500">
                                                                {((chunk.relevance_score || 0) * 100).toFixed(1)}%
                                                            </span>
                                                        </div>
                                                        <p className="text-gray-600 line-clamp-2">
                                                            {chunk.content_preview}
                                                        </p>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>

                                        <ChevronRight className={cn(
                                            "h-4 w-4 transition-transform",
                                            isSelected && "rotate-90"
                                        )} />
                                    </div>
                                </Card>
                            );
                        })}
                    </div>
                </div>
            )}
        </aside>
    );
}