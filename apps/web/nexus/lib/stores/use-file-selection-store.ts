import { create } from 'zustand';

interface FileSelectionState {
    selectedFileIds: Set<string>;
    isSelectionMode: boolean;

    toggleFile: (fileId: string) => void;
    selectedFile: (fileId: string) => void;
    deselectFile: (fileId: string) => void;
    selectAll: (fileId: string[]) => void;
    clearSelection: () => void;
    enterSelectionMode: () => void;
    exitSelectionMode: () => void;

    isSelected: (fileId: string) => boolean;
    getSelectedCount: () => number;
}

export const useFileSelectionStore = create<FileSelectionState>((set, get) => ({
    selectedFileIds: new Set(),
    isSelectionMode: false,

    toggleFile: (fileId: string) => set((state) => {
        const newSet = new Set(state.selectedFileIds);
        if (newSet.has(fileId)) {
            newSet.delete(fileId);
        } else {
            newSet.add(fileId);
        }
        return { selectedFileIds: newSet}
    }),

    selectFile: (fileId: string) => set((state) => {
        const newSet = new Set(state.selectedFileIds);
        newSet.add(fileId);
        return { selectedFileIds: newSet };
    }),

    deselectFile: (fileId: string) => set((state) => {
        const newSet = new Set(state.selectedFileIds);
        newSet.delete(fileId);
        return { selectedFileIds: newSet };
    }),

    selectAll: (fileIds: string[]) => set({
        selectedFileIds: new Set(fileIds),
        isSelectionMode: true,
    }),

    clearSelection: () => set({
        selectedFileIds: new Set(),
        isSelectionMode: false,
    }),

    enterSelectionMode: () => set({ isSelectionMode: true }),

    exitSelectionMode: () => set({ 
        isSelectionMode: false,
        selectedFileIds: new Set(),
    }),

    isSelected: (fileId: string) => {
        return get().selectedFIleIds.has(fileId);
    },

    getSelectedCount: () => {
        return get().selectedFileIds.size;
    }
}));