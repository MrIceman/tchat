package protocol

import (
	"encoding/json"
	"fmt"
	"tchat/internal/message"
)

type ServerSystemMessage struct {
	Message interface{} `json:"Message"`
	msg
}

func NewServerSystemMessage(payload string) *ServerSystemMessage {
	return &ServerSystemMessage{
		Message: payload,
		msg: msg{
			MessageType: message.TypeConnectRes,
		},
	}
}

func (c ServerSystemMessage) Bytes() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(fmt.Errorf("failed to serialize connection request: %s", err.Error()))
	}
	return b
}

func (c ServerSystemMessage) User() string {
	return "server"
}
