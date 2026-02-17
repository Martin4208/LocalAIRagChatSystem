// lib/api/normalizers.ts
import type { ChatMessage } from '../types/chat';

export function normalizeChatMessage(raw: any): ChatMessage {
  return {
    ...raw,
    documentRefs: raw.document_refs,
  };
}
