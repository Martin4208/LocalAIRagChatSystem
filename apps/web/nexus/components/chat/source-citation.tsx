'use client';

// ⚠️ useChatUIStore ではなく useSourcePanelStore を使う
// useChatUIStore には selectSource / selectedSourceIndex が存在しないため
// SourceCitation が undefined になりクラッシュする原因になっていた
import { useSourcePanelStore } from '@/lib/stores/use-source-panel-store';
import { Button } from '@/components/ui/button';

interface SourceCitationProps {
  index: number;
  documentName?: string;
  pageNumber?: number | null;
}

export function SourceCitation({ index, documentName, pageNumber }: SourceCitationProps) {
  const selectSource = useSourcePanelStore((state) => state.selectSource);
  const selectedSourceIndex = useSourcePanelStore((state) => state.selectedSourceIndex);

  const isSelected = selectedSourceIndex === index;

  const label = documentName
    ? pageNumber != null
      ? `${documentName} P.${pageNumber}`
      : documentName
    : `[${index + 1}]`;

  return (
    <Button
      variant="outline"
      size="sm"
      onClick={() => selectSource(index)}
      className={`text-xs px-2 py-1 h-auto max-w-[200px] truncate ${
        isSelected
          ? 'bg-blue-100 border-blue-400 text-blue-700'
          : 'bg-white text-gray-700 hover:bg-gray-100'
      }`}
      title={label}
    >
      {label}
    </Button>
  );
}