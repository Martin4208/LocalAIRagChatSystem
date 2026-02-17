import {
    LucideIcon,
    FileText,
    MessageSquare,
    BarChart3,
    Palette,
    Search,
    Network,
    Upload
} from 'lucide-react';

export interface NavigationItem {
    id: string;
    label: string;
    icon: LucideIcon;
    href: (workspaceId: string) => string;
}

export const NAVIGATION_ITEMS: NavigationItem[] = [
    {
        id: 'dashboard',
        label: 'Dashboard',
        icon: BarChart3,
        href: (workspaceId) => `/workspaces/${workspaceId}`
    },
    // { 
    //     id: 'files',
    //     label: 'Files',
    //     icon: FileText, 
    //     href: (workspaceId) => `/workspaces/${workspaceId}/files` 
    // },
    { 
        id: 'chats',
        label: 'Chats',
        icon: MessageSquare, 
        href: (workspaceId) => `/workspaces/${workspaceId}/chats` 
    },
    // { 
    //     id: 'analysis',
    //     label: 'Analysis',
    //     icon: BarChart3, 
    //     href: (workspaceId) => `/workspaces/${workspaceId}/analysis` 
    // },
    // { 
    //     id: 'canvas',
    //     label: 'Canvas',
    //     icon: Palette, 
    //     href: (workspaceId) => `/workspaces/${workspaceId}/canvas` 
    // },
    // { 
    //     id: 'graph',
    //     label: 'Graph',
    //     icon: Network, 
    //     href: (workspaceId) => `/workspaces/${workspaceId}/graph` 
    // },
    { 
        id: 'upload',
        label: 'Upload',
        icon: Upload, 
        href: (workspaceId) => `/workspaces/${workspaceId}/upload` 
    },
    // { 
    //     id: 'search',
    //     label: 'Search',
    //     icon: Search, 
    //     href: (workspaceId) => `/workspaces/${workspaceId}/search` 
    // },
]