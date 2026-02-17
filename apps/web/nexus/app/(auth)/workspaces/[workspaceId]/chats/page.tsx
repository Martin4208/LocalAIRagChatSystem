'use client';

import Link from 'next/link';
import { useWorkspaceContext } from '@/lib/providers/workspace-context';
import { useChats, useCreateChat, useDeleteChat } from '@/lib/hooks/use-chats';
import { DeleteConfirm } from '@/components/shared/delete-confirm';

export default function ChatsPage() {
    const { workspaceId } = useWorkspaceContext();

    const { data, isLoading, error } = useChats(workspaceId, { limit: 50, offset: 0 });

    const createChatMutation = useCreateChat();

    const deleteChatMutation = useDeleteChat();

    const handleCreateChat = async () => {
        console.log('=== Creating Chat ===');
        console.log('workspaceId:', workspaceId);
        try {
            await createChatMutation.mutateAsync({
                workspaceId, 
                data: {
                    title: `New Chat ${new Date().toLocaleTimeString()}`
                }
            });
            console.log('✅ Chat created successfully:');
        } catch (error) {
            console.log('Failed to create chat:', error);

            // エラーの詳細を表示
            if (error instanceof Error) {
            console.error('Error message:', error.message);
            }
            
            // APIエラーの場合
            if (error?.response) {
                console.error('API Response:', error.response);
                console.error('Status:', error.response.status);
                console.error('Data:', error.response.data);
            }
        }
    }

    const handleDeleteChat = async (chatId: string) => {
        try {
            await deleteChatMutation.mutateAsync({
                workspaceId,
                chatId
            });
        } catch (error) {
            console.log('Failed to delete chat:', error);
        }
    }

    if (isLoading) {
        return <div className="p-8">Loading chats...</div>;
    }

    if (error) {
        return (
        <div className="p-8 text-red-600">
            Error: {error.message}
        </div>
        );
    }

    return (
        <div className="p-8">
            <div className="mb-4 flex items-center justify-between">
                <h1 className="text-2xl font-bold">Chats</h1>
                <button
                    onClick={handleCreateChat}
                    disabled={createChatMutation.isPending}
                    className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                    {createChatMutation.isPending ? 'Creating...' : 'NewChat'} 
                </button>
            </div>

            {/* チャット一覧 */}
            {data?.chats && data.chats.length > 0 ? (
                <div className="space-y-2">
                    {data.chats.map((chat) => (
                        <div 
                            key={chat.id}
                            className="p-4 border rounded flex items-center justify-between"
                        >
                            <Link
                                href={`/workspaces/${workspaceId}/chats/${chat.id}`}
                                className="flex-1 cursor-pointer hover:text-blue-600 transition-colors"
                            >
                                <div>
                                    <h3 className="font-medium">{chat.title}</h3>
                                    <p className="text-sm text-gray-500">
                                        {new Date(chat.created_at).toLocaleString()}
                                    </p>
                                </div>
                                
                            </Link>
                            {/* <button
                                onClick={(e) => {
                                    e.preventDefault();
                                    e.stopPropagation();
                                    handleDeleteChat(chat.id)
                                }}
                                disabled={deleteChatMutation.isPending}
                                className="px-3 py-1 text-red-600 border border-red-600 rounded hover:bg-red-50 disabled:opacity-50"
                            >
                                Delete
                            </button> */}
                            <DeleteConfirm onDelete={() => handleDeleteChat(chat.id)}/>
                        </div>
                    ))}
                </div>
            ): (
                <div className="text-center py-12 text-gray-500">
                    No chats yet. Create your first chat!
                </div>
            )}

            {/* エラー表示 */}
            {createChatMutation.error && (
                <div className="mt-4 p-4 bg-red-50 text-red-600 rounded">
                Failed to create chat: {createChatMutation.error.message}
                </div>
            )}

            {deleteChatMutation.error && (
                <div className="mt-4 p-4 bg-red-50 text-red-600 rounded">
                Failed to delete chat: {deleteChatMutation.error.message}
                </div>
            )}
        </div>
    );
}