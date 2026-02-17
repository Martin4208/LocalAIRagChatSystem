'use client';

import { Button } from '@/components/ui/button';
import { PanelRightOpen, PanelRightClose } from 'lucide-react';
import { useLayoutStore } from '@/lib/stores/use-side-pane-store';

interface ChatHeaderProps {
    title: string;
    messageCount: number;
    hasAnySources: boolean;
}

export function ChatHeader({ title, messageCount, hasAnySources }: ChatHeaderProps) {
    const isSourcePanelOpen = useLayoutStore((state) => state.isSourcePanelOpen);
    const toggleSourcePanel = useLayoutStore((state) => state.toggleSourcePanel);

    return (
        <div className="p-6 flex items-center justify-between">
            <div>
                <h1 className="text-2xl font-bold">{title}</h1>
                <p className="text-sm text-gray-500">
                    {messageCount} messages
                </p>
            </div>

            <Button
                variant="outline"
                size="sm"
                onClick={toggleSourcePanel}
                className="flex items-center gap-2"
                title={isSourcePanelOpen ? "Close Source" : "Open Source"}
            >
                {isSourcePanelOpen ? (
                    <>
                        <PanelRightClose className="h-4 w-4" />
                    </>
                ) : (
                    <>
                        <PanelRightOpen className="h-4 w-4" />
                    </>
                )}
            </Button>
        </div>
    )
}