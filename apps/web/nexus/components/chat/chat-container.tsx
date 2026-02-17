'use client';

import { ReactNode } from 'react';
import { cn } from '@/lib/utils';

interface ChatContainerProps {
  header: ReactNode;
  messages: ReactNode;
  input: ReactNode;
  sources: ReactNode;
}

export function ChatContainer({
    header, 
    messages,
    input,
    sources,
}: ChatContainerProps) {

    return (
        <div 
            className="h-full w-full overflow-hidden flex"
        >
            {/* Chat */}
            <div className="flex flex-col flex-1 min-w-0">

                {/* Header */}
                <div className="flex-none border-b">
                    {header}
                </div>

                {/* Messages (ONLY scroll area) */}
                <div className="flex-1 overflow-y-auto px-4 py-6">
                    {messages}
                </div>

                {/* Input (NOT scrollable) */}
                <div className="flex-none border-t px-4 py-3">
                    {input}
                </div>
            </div>

            {/* Sources */}
            <div
                className={cn(
                    'h-full  bg-gray-50 transition-all duration-300 overflow-hidden'
                )}
            >
                <div className="h-full">{sources}</div>
                
            </div>
        </div>
    );
}