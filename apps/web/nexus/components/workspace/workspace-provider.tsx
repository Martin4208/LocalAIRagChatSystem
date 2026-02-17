interface WorkspaceContextValue {
  workspaceId: string;
  workspace: Workspace | null;
  isLoading: boolean;
}

export const useWorkspace = () => {
  const context = useContext(WorkspaceContext);
  if (!context) throw new Error('useWorkspace must be used within WorkspaceProvider');
  return context;
};