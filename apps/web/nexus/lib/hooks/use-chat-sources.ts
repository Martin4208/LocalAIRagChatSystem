// lib/hooks/use-chat-sources.ts

import { useMemo } from 'react';
import { ChatMessage } from '@/lib/types/chat';
import { 
  extractSourcesFromMessages   
} from '@/lib/utils/chat-sources';

/**
 * チャットメッセージからソース情報を抽出・整形するフック
 */
export function useChatSourcesExtraction(messages: ChatMessage[] | undefined) {
  return useMemo(() => {
    if (!messages || messages.length === 0) {
      return {
        sourcesList: [],
        messageSourcesMap: new Map(),
        hasAnySources: false,
      };
    }

    const { allSources, messageSourcesMap } = extractSourcesFromMessages(messages);

    return {
      sourcesList: allSources,
      messageSourcesMap,
      hasAnySources: allSources.length > 0,
    };
  }, [messages]);
}


/**
 * チャットメッセージから引用ソースを抽出
 */
export function useChatSources(workspaceId: string, chatId: string) {
  const { data: chatData } = useChat(workspaceId, chatId);

  return useMemo(() => {
    // 未使用のため空を返す
    return {
      allSources: [],
      messageSourcesMap: new Map()
    };
  }, []);
}