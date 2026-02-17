'use client'
// 選択されたファイルのプレビュー表示

import { getFileType } from '@/lib/utils/file'
import { DocumentWithFile } from '@/types/domain'
import { useState, useEffect } from 'react';
// import { pdfjs, Document, Page } from 'react-pdf';
// import 'react-pdf/dist/Page/AnnotationLayer.css'
// import 'react-pdf/dist/Page/TextLayer.css'

interface FilePreviewProps {
    document: DocumentWithFile | null
}

// pdfjs.GlobalWorkerOptions.workerSrc = `//unpkg.com/pdfjs-dist@${pdfjs.version}/build/pdf.worker.min.js`;

function PDFPreview({ document }: { document: DocumentWithFile }) {
    // const pdf = document.file.preview_content;
    // const [numPages, setNumPages] = useState<number>(0);
    // const [pageNumber, setPageNumber] = useState<number>(1);
    // const [error, setError] = useState(false);

    // function onDocumentLoadSuccess({ numPages }: {numPages: number }) {
    //     setNumPages(numPages);
    // }

    // if (!pdf) return <div>Loading...</div>;

    // const handleError = () => {
    //     setError(true);
    // }

    // if (error) {
    //     return (
    //         <div className="h-full flex items-center justify-center">
    //             <div className="text-center text-muted-foreground">
    //                 <p className="text-4xl mb-4">⚠️</p>
    //                 <p className="text-lg font-semibold">PDFの読み込みに失敗しました</p>
    //                 <p className="text-sm">{document.name}</p>
    //             </div>
    //         </div>
    //     )
    // }

    // return (
    //     <div className="h-full flex items-center justify-center">
    //         <div className="text-center">
    //             <p className="text-xl font-bold">{document.name}</p>
    //             <Document 
    //                 file={pdf}
    //                 onLoadSuccess={onDocumentLoadSuccess}
    //                 onLoadError={handleError}
    //             >
    //                 <Page pageNumber={pageNumber}/>
    //             </Document>

    //             <div className="flex items-center gap-4">
    //                 <button 
    //                     onClick={() => setPageNumber(prev => prev - 1)}
    //                     disabled={pageNumber <= 1}
    //                 >
    //                     前へ
    //                 </button>
                    
    //                 <span>{pageNumber} / {numPages}</span>
                    
    //                 <button 
    //                     onClick={() => setPageNumber(prev => prev + 1)}
    //                     disabled={pageNumber >= numPages}
    //                 >
    //                     次へ
    //                 </button>
    //             </div>

    //             <p className="text-muted-foreground">PDF Preview</p>
    //         </div>
    //     </div>
    // );
    return (
        <div className="h-full flex items-center justify-center">
            <div className="text-center text-muted-foreground">
                <p className="text-4xl mb-4">📄</p>
                <p className="text-xl font-bold">{document.name}</p>
                <p className="text-lg">PDF Preview</p>
                <p className="text-sm mt-4">
                    PDFプレビュー機能は開発中です
                </p>
                <p className="text-sm text-muted-foreground">
                    デスクトップアプリ版で対応予定
                </p>
            </div>
        </div>
    );
}

function ImagePreview({ document }: { document: DocumentWithFile }) {
    const imageContent = document.file.preview_content;
    const [error, setError] = useState(false);

    if (!imageContent) return <div>No Image</div>;

    const handleError = () => {
        setError(true);
    }

    if (error) {
        return (
            <div className="h-full flex items-center justify-center">
                <div className="text-center text-muted-foreground">
                    <p className="text-4xl mb-4">⚠️</p>
                    <p className="text-lg font-semibold">画像の読み込みに失敗しました</p>
                    <p className="text-sm">{document.name}</p>
                </div>
            </div>
        )
    }
    
    return (
        <div className="h-full flex items-center justify-center p-6 bg-muted/10">
            <div className="text-center">
                <img 
                    src={imageContent} 
                    alt={document.name}
                    onError={() => handleError()}
                    className="max-w-full max-h-full object-contain" 
                />
                <p className="text-muted-foreground">Image Preview</p>
            </div>
        </div>
    );
}

function VideoPreview({ document }: { document: DocumentWithFile }) {
    const video = document.file.preview_content;
    const [error, setError] = useState('');

    if (!video) return <div>Loading...</div>;

    const handleError = () => {
        setError(true);
    }

    if (error) {
        return (
            <div className="h-full flex items-center justify-center">
                <div className="text-center text-muted-foreground">
                    <p className="text-4xl mb-4">⚠️</p>
                    <p className="text-lg font-semibold">ビデオの読み込みに失敗しました</p>
                    <p className="text-sm">{document.name}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="h-full flex items-center justify-center">
            <div className="text-center">
                <p className="text-xl font-bold">{document.name}</p>
                <video 
                    src={video}
                    controls
                    autoPlay={false}
                    loop={false}
                    muted={false}
                    onError={handleError}
                    className="max-w-full max-h-full object-contain"
                ></video>
                <p className="text-muted-foreground">Video Preview</p>
            </div>
        </div>
    );
}

function TextPreview({ document }: { document: DocumentWithFile }) {
    const text = document.file.preview_content;
    const [error, setError] = useState('');

    if (!text) return <div>Loading...</div>;

    const handleError = () => {
        setError(true);
    }

    if (error) {
        return (
            <div className="h-full flex items-center justify-center">
                <div className="text-center text-muted-foreground">
                    <p className="text-4xl mb-4">⚠️</p>
                    <p className="text-lg font-semibold">テキストの読み込みに失敗しました</p>
                    <p className="text-sm">{document.name}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="h-full flex flex-col p-4">
            {/* ヘッダ */}
            <p className="text-xl font-bold mb-2">
                {document.name}
            </p>

            {/* 本文（スクロール領域） */}
            <pre 
                className="flex-1 overflow-auto whitespace-pre-wrap text-sm"
                onError={handleError}
            >
                {text}
            </pre>

            {/* フッタ */}
            <p className="mt-2 text-muted-foreground text-sm">
                Text Preview
            </p>
        </div>
    );
}

function UnsupportedPreview({ document }: { document: DocumentWithFile }) {
    return (
        <div className="h-full flex items-center justify-center">
            <div className="text-center">
                <p className="text-4xl mb-4">❓</p>
                <p className="text-muted-foreground">Unsupported Preview</p>
                <p className="text-sm">Coming soon</p>
            </div>
        </div>
    );
}

export function FilePreview({ document }: FilePreviewProps) {
    if (!document) {
        return (
            <div>
                <p className="text-4xl mb-4">📄</p>
                <p className="text-lg">Select a file to preview</p>
            </div>
        );
    }

    const fileType = getFileType(document);

    switch (fileType) {
        case 'image':
            return <ImagePreview document={document} />
        case 'pdf':
            return <PDFPreview document={document} />
        case 'video':
            return <VideoPreview document={document} />
        case 'text':
            return <TextPreview document={document} />
        default:
            return <UnsupportedPreview document={document} />
    }
}