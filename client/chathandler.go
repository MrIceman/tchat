package client

import "tchat/internal/types"

type ChatHandler struct {
	chatMessages chan types.Message
}

func NewChatHandler(chatMessages chan types.Message) *ChatHandler {
	return &ChatHandler{
		chatMessages: chatMessages,
	}
}
