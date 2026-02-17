// components/chat/chat-input.tsx

'use client';

import { FormEvent, KeyboardEvent, useEffect, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { Send, Loader2, Paperclip } from 'lucide-react';
import { cn } from '@/lib/utils';

interface ChatInputProps {
  value: string;
  onChange: (value: string) => void;
  onSubmit: () => void;
  disabled?: boolean;
  placeholder?: string;
  maxLength?: number;
  showFileButton?: boolean;
  onFileSelect?: (files: FileList) => void;
}

export function ChatInput({
  value,
  onChange,
  onSubmit,
  disabled = false,
  placeholder = 'Message Nexus...',
  maxLength = 2000,
  showFileButton = true,
  onFileSelect,
}: ChatInputProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto-resize logic
  useEffect(() => {
    const textarea = textareaRef.current;
    if (!textarea) return;

    textarea.style.height = '44px'; // Reset to min height
    const scrollHeight = textarea.scrollHeight;
    const newHeight = Math.min(Math.max(scrollHeight, 44), 200);
    textarea.style.height = `${newHeight}px`;
  }, [value]);

  // Form submit handler
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!value.trim() || disabled) return;
    onSubmit();
  };

  // Keyboard shortcuts (Enter: submit, Shift+Enter: new line)
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  // File upload handler
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && onFileSelect) {
      onFileSelect(files);
    }
    e.target.value = ''; // Reset input
  };

  const isOverLimit = value.length > maxLength;
  const canSubmit = value.trim() && !disabled && !isOverLimit;

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-3xl mx-auto">
      <div
        className={cn(
          'flex items-end gap-2 rounded-2xl border bg-background px-3 py-2 shadow-sm transition-all',
          'focus-within:ring-2 focus-within:ring-primary/20 focus-within:border-primary',
          isOverLimit && 'border-destructive focus-within:ring-destructive/20'
        )}
      >
        {/* File upload button */}
        {showFileButton && (
          <div className="relative">
            <input
              type="file"
              id="chat-file-upload"
              multiple
              className="hidden"
              onChange={handleFileChange}
              disabled={disabled}
              aria-label="Upload files"
            />
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="h-9 w-9 rounded-full hover:bg-muted"
              onClick={() => document.getElementById('chat-file-upload')?.click()}
              disabled={disabled}
              aria-label="Attach files"
            >
              <Paperclip className="h-4 w-4" />
            </Button>
          </div>
        )}

        {/* Textarea */}
        <textarea
          ref={textareaRef}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          disabled={disabled}
          maxLength={maxLength}
          className={cn(
            'flex-1 resize-none bg-transparent outline-none border-none',
            'text-sm leading-6 text-foreground placeholder:text-muted-foreground',
            'caret-primary min-h-[44px] max-h-[200px] py-2'
          )}
          aria-label="Chat message input"
          aria-describedby={isOverLimit ? 'char-count-error' : 'char-count'}
        />

        {/* Character count */}
        <div className="flex flex-col items-end gap-1">
          

          {/* Submit button */}
          <Button
            type="submit"
            size="icon"
            className="h-9 w-9 rounded-full"
            disabled={!canSubmit}
            aria-label="Send message"
          >
            {disabled ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Send className="h-4 w-4" />
            )}
          </Button>
        </div>
      </div>

      {/* Error message for screen readers */}
      {isOverLimit && (
        <p className="sr-only" role="alert">
          Message exceeds maximum length of {maxLength} characters
        </p>
      )}
    </form>
  );
}