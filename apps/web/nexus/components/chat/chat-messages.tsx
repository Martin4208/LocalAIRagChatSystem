'use client';

import { useEffect, useRef } from 'react';
import { ChatMessage as ChatMessageType } from '@/lib/types/chat';
import { MessageBubble } from './message-bubble';
import { Button } from '@/components/ui/button';
import { Bot, Loader2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface ChatMessagesProps {
    messages: ChatMessageType[];
    messageSourcesMap: Map<string, number[]>;
    onSuggestedQuestionClick?: (question: string) => void;
    isLoading?: boolean;  
}

const suggestedQuestions = [
    "Summarize the main themes across all documents",
    "What are the key takeaways?",
    "Find connections between the documents",
    "What topics need more exploration?"
];

export function ChatMessages({ 
    messages, 
    messageSourcesMap,
    onSuggestedQuestionClick,
    isLoading
}: ChatMessagesProps) {
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollIntoView({ behavior: 'smooth' });
        } 
    }, [messages]);

    console.log('RAW MESSAGE:', messages);

    if (messages.length === 0) {
        return (
            <motion.div 
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="flex items-center justify-center h-full"
            >
                <div className="text-center py-12">
                    <Bot className="w-16 h-16 mx-auto mb-4 text-foreground-tertiary" />
                    <h3 className="text-lg font-medium text-foreground mb-2">
                        Ask me anything about your documents
                    </h3>
                    <p className="text-sm text-foreground-secondary mb-6 font-mono">
                        I'll search through your knowledge base
                    </p>
                    <div className="flex flex-wrap gap-2 justify-center max-w-md mx-auto">
                        {suggestedQuestions.map((q, i) => (
                            <Button
                                key={i}
                                variant="outline"
                                size="sm"
                                className="text-xs border text-foreground-secondary hover:bg-slate-900 hover:text-primary transition-all"
                                onClick={() => onSuggestedQuestionClick?.(q)}
                            >
                                {q}
                            </Button>
                        ))}
                    </div>
                </div>
            </motion.div>
        );
    }

    return (
        <div className="flex justify-center">
            <div className="w-full max-w-[1120px] px-4">
                <AnimatePresence>
                    <div className="space-y-4 ">
                        {messages.map((message) => (
                            <motion.div
                                key={message.id}
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ duration: 0.3 }}
                            >
                                <MessageBubble
                                    message={message}
                                    sourceIndices={messageSourcesMap.get(message.id) || []}
                                />
                            </motion.div>
                        ))}
                    </div>
                </AnimatePresence>
                
                {/* ローディング表示 */}
                {isLoading && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        className="flex gap-3 mt-4"
                    >
                        <div className="w-8 h-8 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center">
                            <Bot className="w-4 h-4 text-emerald-400" />
                        </div>
                        <div className="bg-background-secondary border rounded-2xl px-4 py-3 shadow-md">
                            <div className="flex items-center gap-2 text-foreground-secondary">
                                <Loader2 className="w-4 h-4 animate-spin text-primary" />
                                <span className="text-sm font-mono">Searching documents...</span>
                            </div>
                        </div>
                    </motion.div>
                )}
                
                <div ref={scrollRef}/>
            </div>
        </div>
    );
}
