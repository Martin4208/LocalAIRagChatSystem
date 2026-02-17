import { create } from 'zustand';

interface WorkspaceState {
    currentWorkspaceId: string | null;
    previousWorkspaceId: string | null;
    setCurrentWorkspace: (id: string) => void;
    clearWorkspace: () => void;
}

export const useWorkspaceStore = create<WorkspaceState>((set) => ({
    currentWorkspaceId: null,
    previousWorkspaceId: null,

    setCurrentWorkspace: (id: string) => set((state) => ({
        currentWorkspaceId: id,
        previousWorkspaceId: state.currentWorkspaceId,
    })),

    clearWorkspace: () => set({
        currentWorkspaceId: null,
        previousWorkspaceId: null,
    }),
}));