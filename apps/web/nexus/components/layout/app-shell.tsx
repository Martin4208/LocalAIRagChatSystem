'use client';

import { ReactNode } from 'react';
import { WorkspaceProvider } from '@/lib/providers/workspace-context';
import { AppHeader } from './app-header';
import { AppSidebar } from './app-sidebar';

interface AppShellProps {
  children: ReactNode;
  workspaceId: string;
}

export function AppShell({ children, workspaceId }: AppShellProps) {
  return (
    <WorkspaceProvider workspaceId={workspaceId}>
      <div className="h-screen w-full overflow-hidden flex">
        {/* Sidebar - Fixed */}
        <AppSidebar />

        <div className="flex flex-col flex-1 min-w-0">
          {/* Header - Fixed */}
          <div className="flex-none">
            <AppHeader />
          </div>
          
          {/* Page content (NO scroll here) */}
          <div className="flex-1 overflow-hidden flex flex-col">
            {children}
          </div>
        </div>
      </div>
    </WorkspaceProvider>
  );
}
