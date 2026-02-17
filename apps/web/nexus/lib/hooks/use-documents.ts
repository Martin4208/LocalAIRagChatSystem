// apps/web/nexus/lib/hooks/use-documents.ts

import { useEffect } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { documents, ProcessDocumentRequest, DocumentStatus } from '@/lib/api/documents';

/**
 * ドキュメント処理を開始するMutation
 */
export function useProcessDocument() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workspaceId,
      documentId,
      options
    }: {
      workspaceId: string;
      documentId: string;
      options?: ProcessDocumentRequest;
    }) => documents.process(workspaceId, documentId, options),

    onSuccess: (_, variables) => {
      // ドキュメント一覧を再取得
      queryClient.invalidateQueries({
        queryKey: ['files', variables.workspaceId]
      });
      
      // ドキュメントステータスを再取得
      queryClient.invalidateQueries({
        queryKey: ['documents', variables.workspaceId, variables.documentId, 'status']
      });
    },

    onError: (error) => {
      console.error('Failed to process document:', error);
    }
  });
}

/**
 * ドキュメントの処理状況を取得（ポーリング対応）
 */
export function useDocumentStatus(
  workspaceId: string,
  documentId: string,
  enabled: boolean = true
) {
  const queryClient = useQueryClient();

  const query = useQuery<DocumentStatus>({
    queryKey: ['documents', workspaceId, documentId, 'status'],
    queryFn: () => documents.getStatus(workspaceId, documentId),
    enabled: enabled && !!workspaceId && !!documentId,
    refetchInterval: (query) => {
      if (query.state.data?.status === 'processing') {
        return 3000;
      }
      return false;
    },
    refetchIntervalInBackground: true,
    retry: 3
  });

  useEffect(() => {
    if (query.data?.status === 'processed' || query.data?.status === 'failed') {
      queryClient.invalidateQueries({
        queryKey: ['files', workspaceId],
        exact: false,
        refetchType: 'active',
      });
    }
  }, [query.data?.status, workspaceId, queryClient]);

  return query;
}

/**
 * ドキュメントのチャンク一覧を取得
 */
export function useDocumentChunks(
  workspaceId: string,
  documentId: string,
  page: number = 1,
  limit: number = 20,
  enabled: boolean = true
) {
  return useQuery({
    queryKey: ['documents', workspaceId, documentId, 'chunks', page, limit],
    queryFn: () => documents.getChunks(workspaceId, documentId, page, limit),
    enabled: enabled && !!workspaceId && !!documentId
  });
}