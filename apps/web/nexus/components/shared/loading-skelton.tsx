import { Skeleton } from '@/components/ui/skeleton'; // shadcn/ui

interface LoadingSkeletonProps {
  variant: 'card' | 'list';
  count?: number;
}

export function LoadingSkeleton({ variant, count = 3 }: LoadingSkeletonProps) {
  if (variant === 'card') {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: count }).map((_, i) => (
          <div key={i} className="rounded-lg border p-4">
            <Skeleton className="h-4 w-3/4 mb-2" />
            <Skeleton className="h-3 w-1/2" />
          </div>
        ))}
      </div>
    );
  }
  
  // list variant
  return (
    <div className="space-y-2">
      {Array.from({ length: count }).map((_, i) => (
        <Skeleton key={i} className="h-12 w-full" />
      ))}
    </div>
  );
}