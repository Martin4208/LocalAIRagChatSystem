import { create } from 'zustand';

interface SourcePanelState {
  // ソースパネルの開閉
  isSourceOpen: boolean;
  
  // 選択されたソース番号（1始まり）
  selectedSourceIndex: number | null;
  
  // アクション
  togglePanel: () => void;
  openPanel: () => void;
  closePanel: () => void;
  selectSource: (index: number) => void;
  clearSelection: () => void;
}

export const useSourcePanelStore = create<SourcePanelState>((set) => ({
    isSourceOpen: false,
    selectedSourceIndex: null,

    togglePanel: () => set((state) => ({
        isSourceOpen: !state.isSourceOpen
    })),

    openPanel: () => set({ isSourceOpen: true }),

    closePanel: () => set({ isSourceOpen: false }),

    selectSource: (index: number) => set({
        selectedSourceIndex: index,
        isSourceOpen: true,
    }),

    clearSelection: () => set({
        selectedSourceIndex: null
    }),
}));