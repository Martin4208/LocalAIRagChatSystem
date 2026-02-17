import { ReactNode } from 'react';
import { cn } from '@/lib/utils';

interface MainContentProps {
  children: ReactNode;
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
  className?: string;
}

export function MainContent({ 
  children, 
  maxWidth = 'xl',
  className 
}: MainContentProps) {
  const maxWidthClass = maxWidth === 'full' 
    ? 'max-w-full' 
    : `max-w-${maxWidth}`;

  return (
    <div className={cn(
      'w-full mx-auto px-6 py-8',
      maxWidthClass,
      className
    )}>
      {children}
    </div>
  );
}