// app/(auth)/workspaces/[workspaceId]/layout.tsx
import { ReactNode } from 'react';
import { AppShell } from '@/components/layout/app-shell';

interface WorkspaceLayoutProps {
  children: ReactNode;
  params: Promise<{ workspaceId: string }>;
}

export default async function WorkspaceLayout({
  children,
  params,  
}: WorkspaceLayoutProps) {
  const { workspaceId } = await params;  
  
  return (
    <AppShell workspaceId={workspaceId}>
      {children}
    </AppShell>
  );
}