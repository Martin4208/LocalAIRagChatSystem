'use client';
// 初期データ取得
import { FilesClient } from './_components/files-client';
import { useWorkspaceContext } from '@/lib/providers/workspace-context';

interface PageProps {
    params: {
        workspaceId: string;
    }
}

export default function FilesPage({ params }: PageProps) {
    const { workspaceId } = useWorkspaceContext();
    // const [documents, directoryTree] = await Promise.all([
    //     getWorkspaceDocuments(workspaceId),
    //     getDirectoryTree(workspaceId)
    // ])

    const directoryTree = [
        {
            id: 'root',
            workspace_id: workspaceId,
            parent_id: null,
            name: 'Root',
            created_at: '2025-01-07T00:00:00Z',
            updated_at: '2025-01-07T00:00:00Z',
            deleted_at: null
        },
    ];

    return (    
        <div>
            <FilesClient 
                workspaceId={workspaceId}
                directoryTree={directoryTree}
            />
        </div>
    )
}