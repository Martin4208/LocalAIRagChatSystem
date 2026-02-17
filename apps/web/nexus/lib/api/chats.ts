import { apiClient } from './client';
import type {
    Chat,
    CreateChatRequest,
    ListChatsResponse,
    GetChatResponse,
    SendMessageRequest,
    SendMessageResponse,
} from '../types/chat';
import { normalizeChatMessage } from './normalizers';

export async function listChats (
    workspaceId: string,
    params?: {
        limit?: number,
        offset?: number,
    }
) :Promise<ListChatsResponse> {
    const queryParams = new URLSearchParams();

    if (params?.limit !== undefined) {
        queryParams.set('limit', params.limit.toString())
    }
    if (params?.offset !== undefined) {
        queryParams.set('offset', params.offset.toString())
    }

    const query = queryParams.toString() 
        ? `?${queryParams.toString()}` 
        : '';

    return apiClient<ListChatsResponse>(
        `/workspaces/${workspaceId}/chats${query}`
    );
}

export async function createChat(
    workspaceId: string,
    data: CreateChatRequest
): Promise<Chat> {
    return apiClient<Chat>(
        `/workspaces/${workspaceId}/chats`,
        {
            method: 'POST',
            body: JSON.stringify(data)
        }
    );
}

export async function getChat(
    workspaceId: string,
    chatId: string,
    params?: {
        message_limit?: number;
    }
): Promise<GetChatResponse> {
    const queryParams = new URLSearchParams();
    if (params?.message_limit) {
        queryParams.set('message_limit', params.message_limit.toString());
    }

    const query = queryParams.toString() ? `?${queryParams.toString()}` : '';

    const response = await apiClient<GetChatResponse>(
        `/workspaces/${workspaceId}/chats/${chatId}${query}`
    );

    return {
        chat: response.chat,
        messages: response.messages.map(normalizeChatMessage),
    }
}

export async function deleteChat(
    workspaceId: string,
    chatId: string
): Promise<void> {
    return apiClient<void>(
        `/workspaces/${workspaceId}/chats/${chatId}`,
        { 
            method: 'DELETE' 
        }
    );
}

export async function sendMessage(
    workspaceId: string,
    chatId: string,
    data: SendMessageRequest
): Promise<SendMessageResponse> {
    const response = await apiClient<SendMessageResponse>(
        `/workspaces/${workspaceId}/chats/${chatId}/messages`,
        {
            method: 'POST',
            body: JSON.stringify(data),
        }
    );

    return {
        user_message: normalizeChatMessage(response.user_message),
        assistant_message: normalizeChatMessage(response.assistant_message),
    }
}