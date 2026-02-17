'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { useWorkspaceContext } from '@/lib/providers/workspace-context';
import { useFiles, useUploadFile, useDeleteFile } from '@/lib/hooks/use-files';
import { useDocumentStatus } from '@/lib/hooks/use-documents';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { 
  Upload, 
  FileText, 
  CheckCircle2, 
  XCircle, 
  Loader2,
  AlertCircle,
  Plus,
  Trash2
} from 'lucide-react';

type UploadedFile = {
  id: string;
  name: string;
  size: number;
  status: 'uploading' | 'processing' | 'processed' | 'failed';
};

export default function UploadPage() {
  const { workspaceId } = useWorkspaceContext();
  const uploadMutation = useUploadFile();

  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const { data: filesData } = useFiles(workspaceId);
  const [isDragging, setIsDragging] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const queryClient = useQueryClient();

  // ファイル選択
  const handleFileSelect = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
    if (event.target) {
        event.target.value = '';
    }
  }, []);

  const handlePlusButtonClick = () => {
    fileInputRef.current?.click();
  };

  useEffect(() => {
    console.log('📊 filesData changed:', filesData);
  }, [filesData]);

  // ドラッグ&ドロップ
  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const file = e.dataTransfer.files[0];
    if (file) {
      setSelectedFile(file);
    }
  }, []);

  // アップロード実行
  const handleUpload = async () => {
    if (!selectedFile) return;

    // ✅ SHA256ハッシュを計算
    const fileHash = await calculateSHA256(selectedFile);

    // ✅ 既存ファイルと照合
    const isDuplicate = filesData?.Files?.some(
      (f: any) => f.SHA256Hash === fileHash
    );

    if (isDuplicate) {
      const confirmed = window.confirm(
      `"${selectedFile.name}" は既にアップロード済みです。\n再度アップロードしますか？`
    );
      
    if (!confirmed) {
        setSelectedFile(null);
        return;
      }
    }

    
    try {
      // アップロード実行
      await uploadMutation.mutateAsync({
        workspaceId,
        file: selectedFile
      });

      // 選択ファイルをクリア
      setSelectedFile(null);
    } catch (error) {
      console.error('Upload failed:', error);
    }
  };

  // ファイルサイズをフォーマット
  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  async function calculateSHA256(file: File): Promise<string> {
    const buffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  }


  return (
    <div className="h-full overflow-y-auto">
      <div className="max-w-4xl mx-auto p-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Upload Documents</h1>
          <p className="text-gray-600">
            Upload PDF or text files to add them to your knowledge base
          </p>
        </div>

        {/* ✅ Hidden File Input */}
        <Input
          ref={fileInputRef}
          type="file"
          className="hidden"
          accept=".pdf,.txt"
          onChange={handleFileSelect}
        />

        {/* Upload Area */}
        <Card className="mb-8">
          <div className="p-6">
            {/* Drag & Drop Zone */}
            <div
              className={`
                border-2 border-dashed rounded-lg p-12 text-center transition-colors
                ${isDragging ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}
                ${selectedFile ? 'bg-gray-50' : ''}
              `}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
            >
              {selectedFile ? (
                // Selected File Preview
                <div className="space-y-4">
                  <FileText className="h-12 w-12 mx-auto text-blue-600" />
                  <div>
                    <p className="font-medium">{selectedFile.name}</p>
                    <p className="text-sm text-gray-500">
                      {formatFileSize(selectedFile.size)}
                    </p>
                  </div>
                  <div className="flex gap-2 justify-center">
                    <Button
                      onClick={handleUpload}
                      disabled={uploadMutation.isPending}
                    >
                      {uploadMutation.isPending ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Uploading...
                        </>
                      ) : (
                        <>
                          <Upload className="h-4 w-4 mr-2" />
                          Upload
                        </>
                      )}
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() => setSelectedFile(null)}
                      disabled={uploadMutation.isPending}
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              ) : (
                // Empty State
                <div className="space-y-4">
                  <Upload className="h-12 w-12 mx-auto text-gray-400" />
                  <div>
                    <p className="text-lg font-medium mb-1">
                      Drop your file here, or{' '}
                      <label className="text-blue-600 hover:text-blue-700 cursor-pointer">
                        browse
                        <Input
                          type="file"
                          className="hidden"
                          accept=".pdf,.txt"
                          onChange={handleFileSelect}
                        />
                      </label>
                    </p>
                    {/* ✅ ＋Button */}
                    <Button
                      size="lg"
                      onClick={handlePlusButtonClick}
                      className="p-0 cursor-pointer"
                      disabled={uploadMutation.isPending}
                    >
                      <Plus className="h-6 w-6" />
                      <p>Browse files</p>
                    </Button>
                    <p className="text-sm text-gray-500">
                      Supports: PDF, TXT (Max 100MB)
                    </p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </Card>

        {/* Uploaded Files List */}
        {filesData?.Files && filesData.Files.length > 0 && (
          <Card>
            <div className="p-6">
              <h2 className="text-xl font-semibold mb-4">Uploaded Files : {(filesData.Total)}</h2>
              <div className="space-y-2">
                {filesData.Files.map((file, index) => {
                  console.log(`File ${index}:`, file)
                  return (
                    <FileStatusCard
                      key={`${file.ID}-${index}`}
                      file={file}
                      workspaceId={workspaceId}
                    />
                  )
                })}
              </div>
            </div>
          </Card>
        )}
      </div>
    </div>
  );
}

// ファイルステータス表示カード
function FileStatusCard({
  file,
  workspaceId,
}: {
  file: any;
  workspaceId: string;
}) {
  // processing中の場合のみポーリング
  const { data: status } = useDocumentStatus(
    workspaceId,
    file.Id,
    file.Status === 'processing'
  );

  const deleteMutation = useDeleteFile();

  const handleFileDelete = async () => {
    const confirmed = window.confirm(
      `${file.FileName}を削除しますか？\n\nこの操作は取り消せません。`
    );

    if (!confirmed) return;

    try {
      await deleteMutation.mutateAsync({
        workspaceId,
        fileId: file.ID,
      });
    } catch (error) {
      console.error('Delete failed:', error);
      alert('削除に失敗しました。もう一度お試しください')
    }
  }

  const getStatusIcon = () => {
    switch (file.Status) {
      case 'uploading':
        return <Loader2 className="h-5 w-5 text-blue-600 animate-spin" />;
      case 'processing':
        return <Loader2 className="h-5 w-5 text-yellow-600 animate-spin" />;
      case 'processed':
        return <CheckCircle2 className="h-5 w-5 text-green-600" />;
      case 'failed':
        return <XCircle className="h-5 w-5 text-red-600" />;
    }
  };

  const getStatusText = () => {
    switch (file.Status) {
      case 'uploading':
        return 'Uploading...';
      case 'processing':
        return status?.progress?.chunks_created 
          ? `Processing... (${status.progress.chunks_created} chunks created)`
          : 'Processing...';
      case 'processed':
        return 'Ready for chat';
      case 'failed':
        return 'Failed';
    }
  };

  const getStatusColor = () => {
    switch (file.status) {
      case 'uploading':
      case 'processing':
        return 'text-yellow-600';
      case 'processed':
        return 'text-green-600';
      case 'failed':
        return 'text-red-600';
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  return (
    <div className="group flex items-center justify-between p-4 border rounded-lg transition-colors">
      <div className="flex items-center gap-3 flex-1">
        {getStatusIcon()}
        <div className="flex-1 min-w-0">
          <p className="font-medium truncate">{file.FileName}</p>
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <span>{formatFileSize(file.SizeBytes)}</span>
            <span>•</span>
            <span className={getStatusColor()}>
              {getStatusText()}
            </span>
          </div>
        </div>
      </div>
      
      <div className="flex items-center gap-2">
        {file.Status === 'processed' && (
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              // チャットページへ遷移
              window.location.href = `/workspaces/${workspaceId}/chats`;
            }}
          >
            Use in Chat
          </Button>
        )}

        <Button
          variant="ghost"
          size="sm"
          onClick={handleFileDelete}
          disabled={deleteMutation.isPending}
          className="
            text-red-600
            hover:bg-red-600
            hover:text-white
            transition-colors
          "
          title="削除"
        >
          {deleteMutation.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Trash2 className="h-4 w-4" />
          )}
        </Button>

        {file.Status === 'failed' && (
          <div className="ml-4">
            <AlertCircle className="h-5 w-5 text-red-600" />
          </div>
        )}
      </div>
    </div>
  );
}