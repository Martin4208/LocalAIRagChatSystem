'use client';

import { useParams, useRouter } from 'next/navigation';
import { useWorkspace } from '@/lib/hooks/use-workspaces';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function WorkspacePage() {
  const params = useParams();
  const router = useRouter();
  const workspaceId = params.workspaceId as string;

  const { data: workspace, isLoading, error } = useWorkspace(workspaceId);

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center min-h-0">
        <p>Loading workspace...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center gap-4 min-h-0">
        <p className="text-red-600">Error: {error.message}</p>
        <Button onClick={() => router.push('/workspaces')}>
          Back to Workspaces
        </Button>
      </div>
    );
  }

  if (!workspace) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center gap-4 min-h-0">
        <p>Workspace not found</p>
        <Button onClick={() => router.push('/workspaces')}>
          Back to Workspaces
        </Button>
      </div>
    );
  }
    // <div className="p-8">
    //   <Card className="p-6">
    //     <h1 className="text-2xl font-bold mb-4">{workspace.name}</h1>
    //     {workspace.description && (
    //       <p className="text-gray-600 mb-4">{workspace.description}</p>
    //     )}
    //     <div className="text-sm text-gray-500">
    //       <p>Created: {new Date(workspace.created_at).toLocaleString()}</p>
    //       <p>Updated: {new Date(workspace.updated_at).toLocaleString()}</p>
    //     </div>
    //   </Card>
    // </div>
    {/* テスト用 */}

  return (

    <div className="space-y-4 p-8">
      <div className="p-4 rounded" style={{ backgroundColor: 'hsl(212, 92%, 45%)', color: 'white' }}>
        Primary色のテスト（インラインスタイル）
      </div>

      <div className="bg-primary text-primary-foreground p-4 rounded">
        Primary色のテスト（Tailwindクラス）
      </div>
      <div className="bg-secondary text-secondary-foreground p-4 rounded">
        Secondary色のテスト
      </div>
      <div className="bg-accent text-accent-foreground p-4 rounded">
        Accent色のテスト
      </div>
      <div className="border-2 border-border p-4 rounded">
        Border色のテスト
      </div>

      <button
        onClick={() => router.push('/workspaces')}
      >
        戻る
      </button>
    </div>  
  );
}