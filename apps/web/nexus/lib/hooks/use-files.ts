// apps/web/nexus/lib/hooks/use-files.ts

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { files } from '@/lib/api/files';

/**
 * ファイル一覧取得
 */
export function useFiles(
  workspaceId: string,
  params?: {
    directoryId?: string;
    limit?: number;
    offset?: number;
  }
) {
  return useQuery({
    queryKey: params 
      ? ['files', workspaceId, params]
      : ['files', workspaceId],
    queryFn: () => files.list(workspaceId, params),
    enabled: !!workspaceId
  });
}

/**
 * ファイルアップロード
 */
export function useUploadFile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workspaceId,
      file,
      directoryId,
      tags
    }: {
      workspaceId: string;
      file: File;
      directoryId?: string;
      tags?: string[];
    }) => files.upload(workspaceId, file, directoryId, tags),

    // アップロード開始時即座にキャッシュ更新
    onMutate: async ({ workspaceId, file }) => {
      console.log('🔵 onMutate START', { workspaceId, fileName: file.name });

      // 進行中のクエリをキャンセル
      await queryClient.cancelQueries({
        queryKey: ['files', workspaceId],
        exact: false,
      });

      const allCaches = queryClient.getQueriesData({
        queryKey: ['files', workspaceId],
        exact: false,
      });

      console.log('All caches: ', allCaches);

      // 現在のキャッシュを保存
      const previousFiles = queryClient.getQueryData(['files', workspaceId]);
      console.log('🔵 previousFiles:', previousFiles);

      // 一時的なファイルオブジェクトを作成
      const tempFile = {
        ID: `temp-${Date.now()}`,
        FileName: file.name,
        SizeBytes: file.size,
        MimeType: file.type,
        Status: 'uploading' as const,
        CreatedAt: new Date().toISOString(),
        SHA256Hash: '',
        Tags: null,
      };

      // キャッシュに即座に追加
      queryClient.setQueryData(
        ['files', workspaceId, undefined],
        (old: any) => {
          console.log('🔵 old cache:', old);

          const newCache = {
            Files: [tempFile, ...(old?.Files || [])],
            total: (old?.total || 0) + 1,
          };
          console.log('🔵 new cache:', newCache);
          return newCache;
        }
      );

      return { previousFiles, tempFile };
    },

    onSuccess: (data, variables, context) => {
      // ファイル一覧を再取得
      queryClient.setQueryData(
        ['files', variables.workspaceId],
        (old: any) => {
          if (!old) return old;

          return {
            ...old,
            Files: old.Files.map((f: any) =>
              f.id === context?.tempFile.id
                ? { ...data, status: 'processing' }
                : f
            ),
          };
        }
      );

      // バックグラウンドで最新データを取得（念のため）
      queryClient.invalidateQueries({
        queryKey: ['files', variables.workspaceId],
        refetchType: 'none', // すぐには再取得しない
      });
    },

    onError: (error, variables, context) => {
      console.error('Failed to upload file:', error);

      if (context?.previousFiles) {
        queryClient.setQueryData(
          ['files', variables.workspaceId],
          context.previousFiles
        );
      }
    },
  });
}

/**
 * ファイル取得
 */
export function useFile(workspaceId: string, fileId: string) {
  return useQuery({
    queryKey: ['files', workspaceId, fileId],
    queryFn: () => files.get(workspaceId, fileId),
    enabled: !!workspaceId && !!fileId
  });
}

/**
 * ファイル削除
 */
export function useDeleteFile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workspaceId,
      fileId,
    }: {
      workspaceId: string;
      fileId: string;
    }) => files.delete(workspaceId, fileId),

    // ✅ Optimistic Update: 即座にUIから削除
    onMutate: async ({ workspaceId, fileId }) => {
      await queryClient.cancelQueries({
        queryKey: ['files', workspaceId],
        exact: false,
      });

      const previousFiles = queryClient.getQueryData([
        'files',
        workspaceId,
        undefined,
      ]);

      // キャッシュから削除
      queryClient.setQueryData(
        ['files', workspaceId, undefined],
        (old: any) => {
          if (!old) return old;

          return {
            ...old,
            Files: old.Files.filter((f: any) => f.ID !== fileId),
            Total: old.Total - 1,
          };
        }
      );

      return { previousFiles };
    },

    // ✅ サーバー削除成功時
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['files', variables.workspaceId],
        exact: false,
      });
    },

    // ✅ エラー時はロールバック
    onError: (error, variables, context) => {
      console.error('❌ Delete failed:', error);

      if (context?.previousFiles) {
        queryClient.setQueryData(
          ['files', variables.workspaceId, undefined],
          context.previousFiles
        );
      }
    },
  });
}

/**
 * ファイルダウンロード
 */
export function useDownloadFile() {
  return useMutation({
    mutationFn: ({
      workspaceId,
      fileId
    }: {
      workspaceId: string;
      fileId: string;
    }) => files.download(workspaceId, fileId)
  });
}