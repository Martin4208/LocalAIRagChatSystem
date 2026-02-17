'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Settings, ChevronLeft, ChevronRight } from 'lucide-react';
import { cn } from '@/lib/utils';
import { LAYOUT } from '@/lib/constants/layout';
import { NAVIGATION_ITEMS } from '@/lib/constants/navigation';
import { useWorkspaceContext } from '@/lib/providers/workspace-context';
import { useLayoutStore } from '@/lib/stores/use-side-pane-store';

export function AppSidebar() {
    const pathname = usePathname();
    const { workspaceId } = useWorkspaceContext();

    const isSidebarOpen = useLayoutStore((state) => state.isSidebarOpen);
    const toggleSidebar = useLayoutStore((state) => state.toggleSidebar);

    const sidebarWidth = isSidebarOpen
        ? LAYOUT.sidebar.width
        : LAYOUT.sidebar.collapsedWidth;

    
    return (
        <aside 
            className="h-full flex flex-col flex-shrink-0 border-r bg-background overflow-hidden
                        transition-[width] duration-300 ease-in-out"
            style={{ width: sidebarWidth }}
        >
            { /* Header */}
            <div
                className="relative flex items-center border-b"
                style={{ height: LAYOUT.header.height }}
            >
                {/* Logo */}
                <div 
                    className={cn(
                        'absolute left-4 text-xl font-bold transition-all duration-300',
                        !isSidebarOpen
                        ? 'opacity-0 -translate-x-2 pointer-events-none'
                        : 'opacity-100 translate-x-0'
                    )}
                >
                    Nexus
                </div>

                {/* Toggle */}
                <button
                    onClick={toggleSidebar}
                    className={cn(
                        'absolute top-4 right-2 z-20',  // ← top位置を調整
                        'w-8 h-8 rounded-full border bg-background shadow-sm',
                        'flex items-center justify-center relative',
                        'transition-all duration-300 hover:bg-accent',
                        !isSidebarOpen
                        ? 'right-2'
                        : 'right-[-16px]'
                    )}
                >
                    {!isSidebarOpen ? (
                        <ChevronRight className="h-5 w-5" />
                    ): (
                        <ChevronLeft className="h-5 w-5" />
                    )}
                </button>
            </div>
            
            {/* Navigation */}
            <nav className="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
                {NAVIGATION_ITEMS.map((item) => {
                    const href = item.href(workspaceId);
                    const isActive = 
                        item.id === 'dashboard'    
                            ? pathname === href
                            : pathname === href || pathname.startsWith(href + '/');
                    const Icon = item.icon;

                    return (
                        <Link
                        key={item.id}
                        href={href}
                        className={cn(
                            'flex items-center h-10 rounded-lg transition-all duration-300',
                            isActive
                            ? 'bg-accent text-accent-foreground'
                            : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                            !isSidebarOpen
                            ? 'justify-center px-2'
                            : 'px-3 gap-3'
                        )}
                        >
                        <Icon className="h-5 w-5 flex-shrink-0" />

                        {/* Label */}
                        <span
                            className={cn(
                            'whitespace-nowrap transition-all duration-300',
                            !isSidebarOpen
                                ? 'opacity-0 max-w-0 overflow-hidden'
                                : 'opacity-100 max-w-[200px]'
                            )}
                        >
                            {item.label}
                        </span>
                        </Link>
                    );
                })}
            </nav>

            {/* Settings */}
            <div className="p-2 border-b">
                <Link
                    href={`/workspaces/${workspaceId}/settings`} 
                    title={!isSidebarOpen ? 'Settings' : undefined}
                    className={cn(
                        'flex items-center h-10 rounded-lg transition-all duration-300 text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                        !isSidebarOpen
                        ? 'justify-center px-2'
                        : 'px-3 gap-3'
                    )}
                >
                    <Settings className="h-5 w-5 flex-shrink-0" />
                    <span
                        className={cn(
                        'transition-all duration-300',
                        !isSidebarOpen
                            ? 'opacity-0 max-w-0 overflow-hidden'
                            : 'opacity-100 max-w-[200px]'
                        )}
                    >
                        Settings
                    </span>
                </Link>
            </div>
        </aside>
    );
}