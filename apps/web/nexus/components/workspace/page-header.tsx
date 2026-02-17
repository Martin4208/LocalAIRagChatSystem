import { ReactNode } from 'react';
import { TYPOGRAPHY } from '@/lib/constants/typography';
import { cn } from '@/lib/utils';

interface PageHeaderProps {
  title: string;
  description?: string;
  action?: ReactNode;
}

export function PageHeader({ title, description, action }: PageHeaderProps) {
  return (
    <div className="flex items-center justify-between mb-6">
      <div className="space-y-1">
        <h1 className={cn(TYPOGRAPHY.heading.h2)}>
          {title}
        </h1>
        {description && (
          <p className={cn(TYPOGRAPHY.body.small)}>
            {description}
          </p>
        )}
      </div>
      {action && (
        <div>
          {action}
        </div>
      )}
    </div>
  );
}