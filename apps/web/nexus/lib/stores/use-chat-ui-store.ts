// lib/stores/use-chat-ui-store.ts
import { create } from 'zustand';

interface ChatUIState {
  // メッセージ入力状態
  draftMessage: string;
  
  // 編集中のメッセージID
  editingMessageId: string | null;
  
  // 返信先メッセージID
  replyToMessageId: string | null;
  
  // スクロール位置を保存すべきか
  shouldPreserveScroll: boolean;
  
  // アクション
  setDraftMessage: (message: string) => void;
  clearDraft: () => void;
  startEditing: (messageId: string) => void;
  stopEditing: () => void;
  setReplyTo: (messageId: string | null) => void;
  setShouldPreserveScroll: (value: boolean) => void;
}

export const useChatUIStore = create<ChatUIState>((set) => ({
  draftMessage: '',
  editingMessageId: null,
  replyToMessageId: null,
  shouldPreserveScroll: false,

  setDraftMessage: (message: string) => set({ draftMessage: message }),
  
  clearDraft: () => set({ draftMessage: '' }),
  
  startEditing: (messageId: string) => set({
    editingMessageId: messageId,
    draftMessage: '', // 編集モードに入ったら下書きをクリア
  }),
  
  stopEditing: () => set({
    editingMessageId: null,
    draftMessage: '',
  }),
  
  setReplyTo: (messageId: string | null) => set({
    replyToMessageId: messageId,
  }),
  
  setShouldPreserveScroll: (value: boolean) => set({
    shouldPreserveScroll: value,
  }),
}));