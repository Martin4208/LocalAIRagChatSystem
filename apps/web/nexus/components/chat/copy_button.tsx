'use client';

import { useState } from 'react';
import { Copy, Check } from 'lucide-react';

interface CopyButtonProps {
    text: string;
}

export function CopyButton({ text }: CopyButtonProps) {
    const [copied, setCopied] = useState(false);

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(text);

            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch {
            console.error('コピーに失敗しました');
        }
    };

    return (
        <button
            onClick={handleCopy}
            className="
                p-1.5
                bg-background/80
                text-foreground-secondary
                hover:text-primary hover:bg-background
                transition
            "
            aria-label="Copy message"
        >
            {copied ? (
                <Check className="w-4 h-4 text-emerald-500"/>
            ) : (
                <Copy className="w-4 h-4" />
            )}
        </button>
    )
}