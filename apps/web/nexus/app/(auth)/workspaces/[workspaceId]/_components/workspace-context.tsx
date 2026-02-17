'use client';

import { createContext, useContext } from 'react';
import { useWorkspace } from '@/lib/hooks/use-workspaces';
import type { Workspace } from '@/types/domain';

type WorkspaceContextValue = {
    workspaceId: string;
    workspace: Workspace = undefined;
    isLoading: boolean;
};

export function WorkspaceProvider({ workspaceId, children }) {
    const { data, isLoading } = useWorkspace(workspaceId);
    
}