// lib/hooks/use-chats.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listChats, createChat, getChat, deleteChat, sendMessage } from '@/lib/api/chats';
import type { CreateChatRequest, SendMessageRequest  } from '@/types/api';

// キャッシュキー
export const chatKeys = {
    all: (workspaceId: string) => 
        ['workspaces', workspaceId, 'chats'] as const,
    detail: (workspaceId: string, chatId: string) => 
        ['workspaces', workspaceId, 'chats', chatId] as const,
};

// 自動キャッシュ、バックグラウンド再フェッチ、ローディング状態管理
export function useChats(
    workspaceId: string,
    params?: { limit?: number, offset?: number }
) {
    return useQuery({
        queryKey: [...chatKeys.all(workspaceId), params],
        queryFn: () => listChats(workspaceId, params),
        enabled: !!workspaceId,
    })
}

export function useChat(
    workspaceId: string, 
    chatId: string,
    params?: { message_limit?: number }
) {
    return useQuery({
        queryKey: [...chatKeys.detail(workspaceId, chatId), params],
        queryFn: () => getChat(workspaceId, chatId, params),
        enabled: !!workspaceId && !!chatId,
    })
}

export function useCreateChat() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ workspaceId, data}: {
            workspaceId: string;
            data: CreateChatRequest
        }) => {
            return createChat(workspaceId, data);
        },
        onSuccess: (newChat, { workspaceId }) => {
            queryClient.invalidateQueries({
                queryKey: chatKeys.all(workspaceId)
            });
        },
    });
}

export function useUpdateChat() {

}

export function useDeleteChat() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ workspaceId, chatId }: { 
            workspaceId: string; 
            chatId: string;
        }) => {
            return deleteChat(workspaceId, chatId);
        },
        onSuccess: (_, { workspaceId, chatId }) => {
            queryClient.invalidateQueries({
                queryKey: chatKeys.all(workspaceId)
            });
            queryClient.invalidateQueries({
                queryKey: chatKeys.detail(workspaceId, chatId)
            });
        },
    });
}

export function useSendMessage() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ workspaceId, chatId, data}: {
            workspaceId: string;
            chatId: string;
            data: SendMessageRequest
        }) => {
            return sendMessage(workspaceId, chatId, data);
        },
        onSuccess: (_, { workspaceId, chatId }) => {
            queryClient.invalidateQueries({ 
                queryKey: chatKeys.detail(workspaceId, chatId) 
            });
        },
    });
}