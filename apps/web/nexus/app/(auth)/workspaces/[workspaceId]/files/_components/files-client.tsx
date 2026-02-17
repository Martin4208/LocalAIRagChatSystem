'use client';

import { useState } from 'react';
import { DirectoryTree } from '@/types/domain';
import { useFiles } from '@/lib/hooks/use-files';
import { FileCard } from './file-card';
import { FileTree } from './file-tree';
import { FilePreview } from './file-preview';
import { Button } from '@/components/ui/button';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Input } from '@/components/ui/input';
import { Search, Upload, Grid, List, Loader2 } from "lucide-react"

interface FilesClientProps {
    workspaceId: string
    directoryTree: DirectoryTree[]
}

type SortBy = 'name' | 'date' | 'size';
type ViewMode = 'tree' | 'grid' | 'list';

export function FilesClient({ 
    workspaceId, 
    directoryTree 
}: FilesClientProps) {
    const { data, isLoading, error } = useFiles(workspaceId);
    const files = data?.files ?? [];

    const [selectedFileId, setSelectedFileId] = useState<string | null>(null);
    const [expandedDirs, setExpandedDirs] = useState<Set<string>>(new Set(['root']))
    const [sortValue, setSortValue] = useState<SortBy>('name');
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [viewMode, setViewMode] = useState<ViewMode>('tree');

    const selectedFile = files.find((file) => file.id === selectedFileId);
    
    const filteredFiles = searchQuery === ''
        ? files
        : files.filter(file => file.name.toLowerCase().includes(searchQuery.toLowerCase()));

    const sortedFiles = [...filteredFiles].sort((a, b) => {
        switch (sortValue) {
            case 'name':
                return a.name.localeCompare(b.name)
            case 'date':
                return new Date(b.created_at).getTime() - new Date(a.created_at).getTime() // 新しい順
            case 'size':
                return b.file.size_bytes - a.file.size_bytes // 大きい順
            default:
                return 0;
        }
    });

    const toggleDirectory = (id: string) => {
        setExpandedDirs(prev => {
            const next = new Set(prev)
            if (next.has(id)) {
                next.delete(id)
            } else {
                next.add(id)
            }
            return next
        })
    }

    const handleFileSelect = (fileId: string) => {
        setSelectedFileId(fileId);
    }

    // ローディング状態
    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-[calc(100vh-73px)]">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    // エラー状態
    if (error) {
        return (
            <div className="flex items-center justify-center h-[calc(100vh-73px)]">
                <p className="text-red-500">Failed to load files: {error.message}</p>
            </div>
        );
    }

    // ステータスに応じたバッジを表示するヘルパー
    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'uploaded':
                return <span className="text-xs px-2 py-1 rounded bg-gray-200 text-gray-700">Uploaded</span>;
            case 'processing':
                return <span className="text-xs px-2 py-1 rounded bg-yellow-200 text-yellow-700">Processing...</span>;
            case 'processed':
                return <span className="text-xs px-2 py-1 rounded bg-green-200 text-green-700">Ready</span>;
            case 'failed':
                return <span className="text-xs px-2 py-1 rounded bg-red-200 text-red-700">Failed</span>;
            default:
                return null;
        }
    };

    return (
        <>
            {/* メインコンテンツ */}
            <div className="flex h-[calc(100vh-73px)]">
                {/* サイドバー（ツリービュー時のみ表示） */}
                {viewMode === 'tree' && (
                    <aside className="w-80 border-r border-border flex flex-col bg-background-secondary">
                        {/* 検索とソート */}
                        <div className="p-4 space-y-3 border-b border-border">
                            {/* 検索 */}
                            <div className="relative">
                                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-foreground-tertiary" />
                                <Input 
                                    placeholder="Search files..."
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    className="pl-9 bg-background-tertiary border-border"
                                />
                                {searchQuery && (
                                    <button
                                        onClick={() => setSearchQuery('')}
                                        className="absolute right-3 top-1/2 -translate-y-1/2 text-foreground-tertiary hover:text-foreground"
                                    >
                                        ✕
                                    </button>
                                )}
                            </div>

                            {/* ソート */}
                            <Select value={sortValue} onValueChange={(value) => setSortValue(value as SortBy)}>
                                <SelectTrigger className="w-full bg-background-tertiary border-border">
                                    <SelectValue placeholder="Sort by..." />
                                </SelectTrigger>
                                <SelectContent className="bg-background-secondary border-border">
                                    <SelectItem value="name">Name</SelectItem>
                                    <SelectItem value="date">Date (Newest)</SelectItem>
                                    <SelectItem value="size">Size (Largest)</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        {/* ファイルツリー */}
                        <div className="flex-1 overflow-y-auto p-4">
                            <p className="text-sm text-muted-foreground">Tree view coming soon...</p>
                        </div>

                        {/* フッター */}
                        <div className="border-t border-border px-4 py-3">
                            <p className="text-sm text-foreground-secondary">
                                {sortedFiles.length} {sortedFiles.length === 1 ? 'file' : 'files'}
                                {searchQuery && ` matching "${searchQuery}"`}
                            </p>
                        </div>
                    </aside>
                )}
                
                {/* メインエリア */}
                <main className="flex-1 overflow-hidden bg-background">
                    {/* ヘッダー */}
                    <div className="p-4 border-b flex items-center justify-between">
                        <div className="flex items-center gap-2">
                            <Button
                                variant={viewMode === 'grid' ? 'default' : 'outline'}
                                size="sm"
                                onClick={() => setViewMode('grid')}
                            >
                                <Grid className="h-4 w-4" />
                            </Button>
                            <Button
                                variant={viewMode === 'list' ? 'default' : 'outline'}
                                size="sm"
                                onClick={() => setViewMode('list')}
                            >
                                <List className="h-4 w-4" />
                            </Button>
                        </div>
                        
                        {/* 検索（グリッド/リストビュー時） */}
                        {viewMode !== 'tree' && (
                            <div className="relative w-64">
                                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                                <Input 
                                    placeholder="Search files..."
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    className="pl-9"
                                />
                            </div>
                        )}
                    </div>

                    {/* ファイル一覧 */}
                    <div className="p-8 overflow-y-auto h-[calc(100%-73px)]">
                        {sortedFiles.length === 0 ? (
                            <div className="flex flex-col items-center justify-center h-full text-center">
                                <p className="text-foreground-secondary">
                                    {searchQuery ? `No files matching "${searchQuery}"` : 'No files yet. Upload some files to get started.'}
                                </p>
                            </div>
                        ) : viewMode === 'grid' ? (
                            // グリッドビュー
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                                {sortedFiles.map((file) => (
                                    <div
                                        key={file.id}
                                        onClick={() => handleFileSelect(file.id)}
                                        className="p-4 rounded-lg border border-border hover:bg-accent/50 transition-colors cursor-pointer"
                                    >
                                        <h3 className="font-medium truncate">{file.fileName}</h3>
                                        <p className="text-sm text-muted-foreground">
                                            {(file.sizeBytes / 1024).toFixed(1)} KB
                                        </p>
                                        <div className="mt-2">
                                            {getStatusBadge(file.status)}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            // リストビュー
                            <div className="space-y-2">
                                {sortedFiles.map((file) => (
                                    <div
                                        key={file.id}
                                        onClick={() => handleFileSelect(file.id)}
                                        className="flex items-center gap-4 p-4 rounded-lg border border-border hover:bg-accent/50 transition-colors cursor-pointer"
                                    >
                                        <div className="flex-1 min-w-0">
                                            <h3 className="font-medium truncate">{file.fileName}</h3>
                                            <p className="text-sm text-muted-foreground">
                                                {new Date(file.createdAt).toLocaleDateString()} · {(file.sizeBytes / 1024).toFixed(1)} KB
                                            </p>
                                        </div>
                                        <div>
                                            {getStatusBadge(file.status)}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </main>
            </div>
        </>
    )
}
