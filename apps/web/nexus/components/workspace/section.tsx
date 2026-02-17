import { ReactNode } from 'react';
import { cn } from '@/lib/utils';

interface SectionProps {
  title?: string;
  children: ReactNode;
  className?: string;
}

export function Section({ title, children, className }: SectionProps) {
  return (
    <section className={cn('space-y-4', className)}>
      {title && (
        <h2 className="text-lg font-semibold">
          {title}
        </h2>
      )}
      {children}
    </section>
  );
}
