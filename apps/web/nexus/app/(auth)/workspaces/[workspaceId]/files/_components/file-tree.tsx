// ツリー操作
'use client'

import { DirectoryTree } from '@/types/domain'
import { cn } from "@/lib/utils"

interface FileTreeProps {
    tree: DirectoryTree[];
    documents: DocumentWithFile[];
    selectedFileId: string | null;
    expandedDirs: Set<string>;
    onToggleDirectory: (id: string) => void;
    onFileSelect: (id: string) => void;
}

export function FileTree({ tree, documents, selectedFileId, expandedDirs, onToggleDirectory, onFileSelect}: FileTreeProps) {

    return (
        <div className="p-4">
            <h2 className="font-bold mb-4">Directories</h2>
            <ul>
                {tree.map((dir) => {
                    const isExpanded = expandedDirs.has(dir.id)
                    
                    return (
                        <li key={dir.id}>
                            <button onClick={() => onToggleDirectory(dir.id)}>
                                {isExpanded ? '📂' : '📁'} {dir.name}
                            </button>

                            {isExpanded && (
                                <div className="pl-4">
                                    {documents
                                        .filter(doc => doc.directory_id === dir.id)
                                        .map(doc => (
                                            <button 
                                                key={doc.id}
                                                onClick={() => onFileSelect(doc.id)}
                                            >
                                                📄 {doc.name}
                                            </button>
                                        ))
                                    }
                                </div>
                            )}
                        </li>
                    )
                })}
            </ul>
                
        </div>
    )
}



