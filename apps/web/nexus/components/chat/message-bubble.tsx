// components/chat/message-bubble.tsx

'use client';

import { ChatMessage } from '@/lib/types/chat';
import { SourceCitation } from './source-citation';
import { Bot, User, FileText } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { CopyButton } from './copy_button';

interface MessageBubbleProps {
    message: ChatMessage;
    sourceIndices: number[];
}

export function MessageBubble({ message, sourceIndices }: MessageBubbleProps) {
    const isUser = message.role === 'user';

    return (
        <div className={`flex gap-2 ${isUser ? 'justify-end' : 'justify-start'}`}>
            {!isUser && (
                <div className="w-8 h-8 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center flex-shrink-0">
                    <Bot className="w-4 h-4 text-primary" />
                </div>
            )}
            <div className={`flex flex-col gap-1 ${isUser ? 'items-end' : 'items-start'}`}>
                <div
                    className={`relative max-w-[70%] rounded-lg p-4 shadow-lg ${
                    isUser
                        ? 'bg-primary text-primary-foreground'
                        : 'bg-background-secondary text-foreground'
                    }`}
                >
                    {!isUser ? (
                        <div className="prose prose-sm prose-invert max-w-none">
                        <ReactMarkdown>{message.content}</ReactMarkdown>
                        </div>
                    ) : (
                        <p className="text-sm">{message.content}</p>
                    )}

                    {/* Source citations - AI messages only */}
                    {!isUser && sourceIndices.length > 0 && (
                        <div className="mt-2 pt-2 border-t border">
                            <p className="text-xs text-foreground-secondary mb-1 flex items-center gap-1">
                                <FileText className="w-3 h-3" /> Sources:
                            </p>
                            <div className="flex flex-wrap gap-1">
                                {sourceIndices.map((index) => (
                                    <SourceCitation key={index} index={index} />
                                ))}
                            </div>
                        </div>
                    )}
                </div>
                
                <CopyButton text={message.content} />
            </div>
            
            {isUser && (
                <div className="w-8 h-8 rounded-full bg-background-secondary border border flex items-center justify-center flex-shrink-0">
                    <User className="w-4 h-4 text-foreground-secondary" />
                </div>
            )}
        </div>
    );
}