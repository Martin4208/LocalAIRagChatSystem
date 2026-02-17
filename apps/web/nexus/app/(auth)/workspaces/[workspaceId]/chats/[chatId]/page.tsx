'use client';

import { useState, useRef, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useChat, useSendMessage } from '@/lib/hooks/use-chats';
import { useChatSourcesExtraction } from '@/lib/hooks/use-chat-sources';

import { ChatContainer } from '@/components/chat/chat-container';
import { ChatHeader } from '@/components/chat/chat-header';
import { ChatMessages } from '@/components/chat/chat-messages';
import { ChatInput } from '@/components/chat/chat-input';
import { SourcesPanel } from '@/components/chat/sources-panel';

export default function Chat() {
    const params = useParams();
    const [input, setInput] = useState('');

    const workspaceId = params.workspaceId === '123'
                ? '550e8400-e29b-41d4-a716-446655440000'
                : params.workspaceId;
    const chatId = params.chatId === '123'
                ? '660e8400-e29b-41d4-a716-446655440001'
                : params.chatId;

    

    const { data, isLoading, error } = useChat(workspaceId, chatId);
    const sendMessageMutation = useSendMessage();

    console.log("=== DEBUG ===");
    console.log("data:", data);
    console.log("data.messages:", data?.messages);
    console.log("data.chat:", data?.chat);

    const { sourcesList, messageSourcesMap, hasAnySources } = useChatSourcesExtraction(
        data?.messages
    );

    const handleSend = async () => {
        if (!input.trim() || sendMessageMutation.isPending) return;

        try {
        await sendMessageMutation.mutateAsync({
            workspaceId,
            chatId,
            data: {
            content: input,
            },
        });
        setInput('');
        } catch (error) {
        console.error('Chat error:', error);
        }
    };

    // ローディング状態
    if (isLoading) {
        return (
        <div className="flex items-center justify-center h-screen">
            <p>Loading chat...</p>
        </div>
        );
    }

    // エラー状態
    if (error) {
        return (
        <div className="flex items-center justify-center h-screen">
            <p className="text-red-600">Error: {error.message}</p>
        </div>
        );
    }

    // データなし
    if (!data) {
        return null;
    }

    return (
    <ChatContainer
        header={
            <ChatHeader
                title={data.chat.title}
                messageCount={data.messages.length}
                hasAnySources={hasAnySources}
            />
        }
        messages={
            <ChatMessages
                messages={data.messages}
                messageSourcesMap={messageSourcesMap}
                onSuggestedQuestionClick={setInput} 
                isLoading={sendMessageMutation.isPending}  
            />
        }
        input={
            <ChatInput
                value={input}
                onChange={setInput}
                onSubmit={handleSend}
                disabled={sendMessageMutation.isPending}
            />
        }
        sources={
            <SourcesPanel
                workspaceId={workspaceId}
                chatId={chatId}
                sourcesList={sourcesList}
            />
        }
        />
    );
}
