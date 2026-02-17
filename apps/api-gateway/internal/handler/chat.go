package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sqlc-dev/pqtype"
)

// ========================================
// CreateChat
// ========================================

func (h *Handler) CreateChat(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
) {
	ctx := r.Context()

	var reqBody api.CreateChatJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if reqBody.Title == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Title is required")
		return
	}

	filterConfig := pqtype.NullRawMessage{
		Valid: false,
	}

	if reqBody.FilterConfig != nil {
		configJSON, err := json.Marshal(reqBody.FilterConfig)
		if err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid filter_config")
			return
		}
		filterConfig = pqtype.NullRawMessage{
			RawMessage: configJSON,
			Valid:      true,
		}
	}

	chat, err := h.queries.CreateChat(ctx, db.CreateChatParams{
		WorkspaceID:  workspaceId,
		Title:        reqBody.Title,
		FilterConfig: filterConfig,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to create chat")
		return
	}

	response := chatToAPI(chat, 0, nil)
	respondJSON(w, http.StatusCreated, response)
}

// ========================================
// ListChats
// ========================================

func (h *Handler) ListChats(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	params api.ListChatsParams,
) {
	ctx := r.Context()

	limit := int32(50)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}

	offset := int32(0)
	if params.Offset != nil {
		offset = int32(*params.Offset)
	}

	chats, err := h.queries.ListChats(ctx, db.ListChatsParams{
		WorkspaceID: workspaceId,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to list chats")
		return
	}

	total, err := h.queries.CountChats(ctx, workspaceId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to count chats")
		return
	}

	apiChats := make([]api.Chat, len(chats))
	for i, chat := range chats {
		// メッセージ数を取得
		count, _ := h.queries.GetChatMessageCount(ctx, chat.ID)

		// 最終メッセージ時刻を取得
		var lastMessageAt *time.Time
		if lastTime, err := h.queries.GetLastMessageTime(ctx, chat.ID); err == nil {
			// swlc が sql.NullTime を返す場合
			if !lastTime.IsZero() {
				t := lastTime
				lastMessageAt = &t
			}
		}

		apiChats[i] = listChatRowToAPI(chat, int(count), lastMessageAt)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"chats": apiChats,
		"total": total,
	})
}

// ========================================
// GetChat
// ========================================

func (h *Handler) GetChat(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	chatId openapi_types.UUID,
	params api.GetChatParams,
) {
	ctx := r.Context()

	chat, err := h.queries.GetChat(ctx, db.GetChatParams{
		ID:          chatId,
		WorkspaceID: workspaceId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Chat not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to get chat")
		return
	}

	messageLimit := int32(50)
	if params.MessageLimit != nil {
		messageLimit = int32(*params.MessageLimit)
	}

	messages, err := h.queries.GetChatMessages(ctx, db.GetChatMessagesParams{
		ChatID: chatId,
		Limit:  messageLimit,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to get messages")
		return
	}

	// 最終メッセージ時刻を取得
	var lastMessage *db.ChatMessage
	if len(messages) > 0 {
		// メッセージがある場合、最後のメッセージを使用
		lastMessage = &messages[len(messages)-1]
	}

	apiChat := chatToAPI(chat, int64(len(messages)), lastMessage)
	apiMessages := make([]api.ChatMessage, len(messages))
	for i, msg := range messages {
		apiMessages[i] = messageToAPI(msg)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"chat":     apiChat,
		"messages": apiMessages,
	})
}

// ========================================
// DeleteChat
// ========================================

func (h *Handler) DeleteChat(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	chatId openapi_types.UUID,
) {
	ctx := r.Context()

	err := h.queries.DeleteChat(ctx, db.DeleteChatParams{
		ID:          chatId,
		WorkspaceID: workspaceId,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to delete chat")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ========================================
// SendMessage
// ========================================

func (h *Handler) SendMessage(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	chatId openapi_types.UUID,
) {
	ctx := r.Context()

	var reqBody api.SendMessageJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if reqBody.Content == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Content is required")
		return
	}

	_, err := h.queries.GetChat(ctx, db.GetChatParams{
		ID:          chatId,
		WorkspaceID: workspaceId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "CHAT_NOT_FOUND", "Chat not found")
			return
		} else {
			log.Printf("Failed to get chat: %v", err)
			respondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch chat")
		}
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	qtx := h.queries.WithTx(tx)

	maxIndex, err := qtx.GetMaxMessageIndex(ctx, chatId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to get message index")
		return
	}

	userMessage, err := qtx.CreateChatMessage(ctx, db.CreateChatMessageParams{
		ChatID:       chatId,
		Role:         "user",
		Content:      reqBody.Content,
		MessageIndex: maxIndex + 1,
		DocumentRefs: pqtype.NullRawMessage{Valid: false},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to save user message")
		return
	}

	assistantContent, documentRefs, err := h.chatService.GenerateResponse(ctx, workspaceId, chatId, reqBody.Content)
	log.Printf("🧪 [Handler] documentRefs len=%d value=%+v", len(documentRefs), documentRefs)
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		// エラー時はフォールバックメッセージ
		assistantContent = "申し訳ございません。応答の生成中にエラーが発生しました。システム管理者に連絡してください。"
	}

	var docRefs pqtype.NullRawMessage

	if len(documentRefs) > 0 {
		b, err := json.Marshal(documentRefs)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "SERIALIZE_ERROR", "Failed to serialize document refs")
			return
		}

		docRefs = pqtype.NullRawMessage{
			RawMessage: b,
			Valid:      true,
		}
	} else {
		docRefs = pqtype.NullRawMessage{Valid: false}
	}

	assistantMessage, err := qtx.CreateChatMessage(ctx, db.CreateChatMessageParams{
		ChatID:       chatId,
		Role:         "assistant",
		Content:      assistantContent,
		MessageIndex: maxIndex + 2,
		DocumentRefs: docRefs,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to save assistant message")
		return
	}

	if err := qtx.UpdateChatTimestamp(ctx, chatId); err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to update chat timestamp")
		return
	}

	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"user_message":      messageToAPI(userMessage),
		"assistant_message": messageToAPI(assistantMessage),
	})
}

// ========================================
// Helper Functions
// ========================================

func chatToAPI(chat db.Chat, messageCount int64, lastMessage *db.ChatMessage) api.Chat {
	var filterConfig *api.FilterConfig
	if chat.FilterConfig.Valid && len(chat.FilterConfig.RawMessage) > 0 {
		var config api.FilterConfig
		if err := json.Unmarshal(chat.FilterConfig.RawMessage, &config); err == nil {
			filterConfig = &config
		}
	}

	var msgCount *int
	if messageCount > 0 {
		count := int(messageCount)
		msgCount = &count
	}

	var lastMessageAt *time.Time
	if lastMessage != nil {
		lastMessageAt = &lastMessage.CreatedAt
	}

	return api.Chat{
		Id:            chat.ID,
		WorkspaceId:   chat.WorkspaceID,
		Title:         chat.Title,
		FilterConfig:  filterConfig,
		MessageCount:  msgCount,
		LastMessageAt: lastMessageAt,
		CreatedAt:     chat.CreatedAt,
		UpdatedAt:     chat.UpdatedAt,
	}
}

func listChatRowToAPI(chat db.ListChatsRow, messageCount int, lastMessageAt *time.Time) api.Chat {
	var filterConfig *api.FilterConfig
	if chat.FilterConfig.Valid {
		var config api.FilterConfig
		json.Unmarshal(chat.FilterConfig.RawMessage, &config)
		filterConfig = &config
	}

	var msgCount *int
	if messageCount > 0 {
		msgCount = &messageCount
	}

	return api.Chat{
		Id:            chat.ID,
		WorkspaceId:   chat.WorkspaceID,
		Title:         chat.Title,
		FilterConfig:  filterConfig,
		MessageCount:  msgCount,
		LastMessageAt: lastMessageAt,
		CreatedAt:     chat.CreatedAt,
		UpdatedAt:     chat.UpdatedAt,
	}
}

func messageToAPI(msg db.ChatMessage) api.ChatMessage {
	var docRefs *[]api.DocumentReference
	if msg.DocumentRefs.Valid {
		var refs []api.DocumentReference
		json.Unmarshal(msg.DocumentRefs.RawMessage, &refs)
		docRefs = &refs
	}

	return api.ChatMessage{
		Id:           msg.ID,
		ChatId:       msg.ChatID,
		Role:         api.ChatMessageRole(msg.Role),
		Content:      msg.Content,
		MessageIndex: int(msg.MessageIndex),
		DocumentRefs: docRefs,
		CreatedAt:    msg.CreatedAt,
	}
}
