import { useQuery } from '@tanstack/react-query';
import { sourcesApi } from '@/lib/api/sources';

export function useChatSources(workspaceId: string, chatId: string) {
  return useQuery({
    queryKey: ['sources', 'chat', workspaceId, chatId],
    queryFn: () => sourcesApi.getChatSources(workspaceId, chatId),
    enabled: !!workspaceId && !!chatId,
  });
}

export function useAnalysisSources(workspaceId: string, analysisId: string) {
  return useQuery({
    queryKey: ['sources', 'analysis', workspaceId, analysisId],
    queryFn: () => sourcesApi.getAnalysisSources(workspaceId, analysisId),
    enabled: !!workspaceId && !!analysisId,
  });
}