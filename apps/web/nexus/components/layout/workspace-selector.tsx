'use client';

import { useParams, useRouter } from 'next/navigation';
import { ChevronDown, Plus, Check } from 'lucide-react';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useWorkspaces } from '@/lib/hooks/use-workspaces';

export function WorkspaceSelector() {
    const params = useParams();
    const router = useRouter();
    const workspaceId = params.workspaceId as string;

    const { data, isLoading } = useWorkspaces();
    const currentWorkspace = data?.workspaces.find((w) => w.id === workspaceId);

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <button className="w-full text-left px-3 py-2 bg-background-tertiary hover:bg-background-tertiary/80 rounded-lg flex items-center justify-between transition-colors">
                    {isLoading 
                    ? <span>Loading...</span>
                    : <span>{currentWorkspace?.name || 'Select Workspace'}</span>
                    }
                    <ChevronDown className="h-4 w-4"/>
                </button>
            </DropdownMenuTrigger>

            <DropdownMenuContent className="w-64 bg-background-secondary border-border">
                    {data?.workspaces.map((workspace) => (
                        <DropdownMenuItem
                            key={workspace.id}
                            onClick={() => router.push(`/workspaces/${workspace.id}`)}
                            className="flex items-center justify-between cursor-pointer"
                        >
                            <span className="text-sm">{workspace.name}</span>
                            {workspace.id === workspaceId && (
                                <Check className="h-4 w-4 text-primary" />
                            )}
                        </DropdownMenuItem>
                    ))}

                    <DropdownMenuSeparator />

                    <DropdownMenuItem
                        onClick={() => router.push(`/workspaces/new`)}
                        className="flex items-center gap-2 text-primary"
                    >
                        <Plus className="h-4 w-4"/>
                        <span>New Workspace</span>
                    </DropdownMenuItem>
                <DropdownMenuSeparator />
            </DropdownMenuContent>
        </DropdownMenu>
    );
}