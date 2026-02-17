'use client';

import { usePathname } from 'next/navigation';
import Link from 'next/link';
import { cn } from '@/lib/utils';

export function TabNavigation({ workspaceId }: { workspaceId: string }) {
    const pathname = usePathname();

    const tabs = [
        { name: 'Files', href: `/workspaces/${workspaceId}/files` },
        { name: 'Chat', href: `/workspaces/${workspaceId}/chat` },
        { name: 'Analysis', href: `/workspaces/${workspaceId}/analysis` },
        { name: 'Canvas', href: `/workspaces/${workspaceId}/canvas` },
        { name: 'Graph', href: `/workspaces/${workspaceId}/graph` },
        { name: 'Upload', href: `/workspaces/${workspaceId}/upload` },
        { name: 'Search', href: `/workspaces/${workspaceId}/search` },
    ];

    return (
        <nav>
            {tabs.map((tab) => {
                const isActive = pathname.startsWith(tabHref);
                
                return (
                    <Link
                        key={tab.href}
                        href={tab.href}
                        className={cn(
                            "px-4 py-2 text-sm font-medium transition-colors",
                            isActive
                                ? "border-b-2 border-primary text-primary"  // アクティブ
                                : "text-muted-foreground hover:text-foreground"  // 非アクティブ
                        )}
                    >
                        {tab.name}
                    </Link>
                );
            })}
        </nav>
    );
}