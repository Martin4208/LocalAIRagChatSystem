// lib/stores/use-side-pane-store.ts
import { create } from 'zustand';

interface LayoutStore {
  isSidebarOpen: boolean;
  isSourcePanelOpen: boolean;
  toggleSidebar: () => void;
  toggleSourcePanel: () => void;
}

export const useLayoutStore = create<LayoutStore>((set) => ({
  isSidebarOpen: true,
  isSourcePanelOpen: true,
toggleSidebar: () =>
    set((s) => ({ isSidebarOpen: !s.isSidebarOpen })),

  toggleSourcePanel: () =>
    set((s) => ({ isSourcePanelOpen: !s.isSourcePanelOpen })),
}));
