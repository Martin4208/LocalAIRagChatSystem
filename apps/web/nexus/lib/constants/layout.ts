export const LAYOUT = {
    sidebar: {
        width: 200,
        widthClass: 'w-64',
        collapsedWidth: 64,
        collapsedWidthClass: 'w-16',
    },
    header: {
        height: 64,
        heightClass: 'h-16',
    },
    source: {
        width: {
            closed: 0,
            open: 300,
        }
    },
    content: {
        maxWidth: {
            sm: 640,
            md: 768,
            lg: 1024,
            xl: 1280,
            full: '100%',
        },
        padding: {
            x: 24,
            y: 32,
        },
    },
    breakpoints: {
        mobile: 768,
    },
} as const