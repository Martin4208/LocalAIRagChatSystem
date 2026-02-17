// apps/web/nexus/app/(auth)/workspaces/[workspaceId]/files/_components/file-card.tsx

'use client';

import { Card } from '@/components/ui/card';
import { Loader2, CheckCircle2, XCircle } from 'lucide-react';
import { useDocumentStatus } from '@/lib/hooks/use-documents';
import { useParams } from 'next/navigation';
import type { DocumentWithFile } from '@/types/domain';

interface FileCardProps {
    document: DocumentWithFile;
}

export function FileCard({ document }: FileCardProps) {
    const params = useParams();
    const workspaceId = params.workspaceId as string;

    // ドキュメントの処理状況を取得（自動ポーリング）
    const { data: status } = useDocumentStatus(
        workspaceId,
        document.id,
        true // 常に有効（内部でstatusに応じてポーリング停止）
    );

    const formatBytes = (bytes: number) => {
        if (bytes < 1024) return `${bytes} B`;
        if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
        return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
    };

    return (
        <Card className="h-[300px] overflow-hidden">
            <div className="h-[200px] bg-muted flex items-center justify-center">
                <p className="text-muted-foreground">Preview</p>
            </div>

            <div className="h-[100px] p-4">
                <p className="font-medium truncate">{document.name}</p>
                <p className="text-sm text-muted-foreground">
                    {formatBytes(document.file.size_bytes)}
                </p>
                <p className="text-sm text-muted-foreground">
                    {new Date(document.created_at).toLocaleDateString()}
                </p>

                {/* 処理ステータス表示 */}
                <div className="mt-2">
                    {status?.status === 'processing' && (
                        <div className="flex items-center gap-2 text-sm text-blue-600">
                            <Loader2 className="h-4 w-4 animate-spin" />
                            <span>Processing... {status.progress?.percentage || 0}%</span>
                        </div>
                    )}

                    {status?.status === 'processed' && (
                        <div className="flex items-center gap-2 text-sm text-green-600">
                            <CheckCircle2 className="h-4 w-4" />
                            <span>Ready</span>
                        </div>
                    )}

                    {status?.status === 'failed' && (
                        <div className="flex items-center gap-2 text-sm text-red-600">
                            <XCircle className="h-4 w-4" />
                            <span>Failed</span>
                        </div>
                    )}
                </div>
            </div>
        </Card>
    );
}