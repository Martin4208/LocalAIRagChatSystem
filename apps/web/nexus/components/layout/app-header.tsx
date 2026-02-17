'use client';

import { useWorkspaceContext } from '@/lib/providers/workspace-context';
import { LAYOUT } from '@/lib/constants/layout';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { useWorkspace } from '@/lib/hooks/use-workspaces';

export function AppHeader() {
    const { workspaceId } = useWorkspaceContext();
    const { data: workspace, isLoading } = useWorkspace(workspaceId);

    return (
        <header 
            className="w-full border-b bg-background"
            style={{ 
                height: LAYOUT.header.height,
                // left: LAYOUT.sidebar.width,
                // width: `calc(100% - ${LAYOUT.sidebar.width}px)`,
            }}
        >
            <div className="flex items-center justify-between h-full px-6">
                {/* Workspace Name */}
                <div className="flex items-center gap-3">
                    {isLoading ? (
                        <LoadingSpinner />
                    ) : (
                        <h2 className="text-lg font-semibold">
                            {workspace?.name || 'Workspace'}
                        </h2>
                    )}
                </div>
            </div>
        </header>
    );
}