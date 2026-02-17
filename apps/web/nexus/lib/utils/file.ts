import type { DocumentWithFile } from '@/types/domain'

export type FileType = 'image' | 'pdf' | 'video' | 'text' | 'unknown'

const IMAGE_EXTENSIONS = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp']
const VIDEO_EXTENSIONS = ['mp4', 'webm', 'mov', 'avi', 'mkv']
const PDF_EXTENSIONS = ['pdf']
const TEXT_EXTENSIONS = ['txt', 'md', 'json', 'csv', 'log']

export function getFileExtension(filename: string): string {
    return filename.split('.').pop()?.toLowerCase() ?? "";
}

export function getFileType(document: DocumentWithFile): FileType {
    const mimeType = document.file.mime_type

    // 画像
    if (mimeType.startsWith('image/')) return 'image';

    // PDF
    if (mimeType === 'application/pdf') return 'pdf';

    // 動画
    if (mimeType.startsWith('video/')) return 'video';

    // テキスト
    if (mimeType.startsWith('text/')) return 'text';

    // ファイル名を取得
    const filename = document.file.original_filename || document.name

    // 拡張子を抽出
    const ext = getFileExtension(filename)

    // 拡張子リストと照合
    if (IMAGE_EXTENSIONS.includes(ext)) return 'image'
    if (PDF_EXTENSIONS.includes(ext)) return 'pdf'
    if (VIDEO_EXTENSIONS.includes(ext)) return 'video'
    if (TEXT_EXTENSIONS.includes(ext)) return 'text'

    return 'unknown'
}