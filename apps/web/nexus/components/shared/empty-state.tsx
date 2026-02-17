import { cn } from '@/lib/utils';

interface EmptyStateProps {
  icon?: React.ReactNode;      // アイコン（lucide-react）
  title: string;                // 見出し
  description?: string;         // 説明文
  action?: React.ReactNode;     // アクションボタン
}

export function EmptyState({ icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      {/* アイコン */}
      {icon && (
        <div className="mb-4">
          {icon}
        </div>
      )}
      
      {/* タイトル */}
      <h3 className="text-lg font-semibold">
        {title}
      </h3>
      
      {/* 説明 */}
      {description && (
        <p className="mt-2 text-sm text-muted-foreground max-w-sm">
          {description}
        </p>
      )}
      
      {/* アクション */}
      {action && (
        <div className="mt-4">
          {action}
        </div>
      )}
    </div>
  );
}