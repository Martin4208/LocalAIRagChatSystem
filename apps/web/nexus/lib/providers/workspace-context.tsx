'use client';

import { createContext, useContext, ReactNode } from 'react';

interface WorkspaceContextValue {
  workspaceId: string;
}

interface WorkspaceProviderProps {
  children: ReactNode;
  workspaceId: string;
}

const WorkspaceContext = createContext<WorkspaceContextValue | undefined>(undefined);

export function useWorkspaceContext() {
  const context = useContext(WorkspaceContext);
    
  if (context === undefined) {
    throw new Error('useWorkspaceContext must be used within WorkspaceProvider');
  }
  
  return context;
}

export function WorkspaceProvider({ children, workspaceId }: WorkspaceProviderProps) {  
  return (
    <WorkspaceContext.Provider value={{ workspaceId }}>
      {children}
    </WorkspaceContext.Provider>
  );
}