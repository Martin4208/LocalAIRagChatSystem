export type FileWithStatus = FileMetadata & {
    status?: 'uploading' | 'processing' | 'processed' | 'failed';
};
