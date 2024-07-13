package protocol

import (
	"encoding/json"
	"fmt"
	"tchat/internal/message"
)

type ClientConnectMessage struct {
	Payload []byte `json:"payload"`

	msg
}

func NewClientConnectMessage(userID string) *ClientConnectMessage {
	return &ClientConnectMessage{
		Payload: []byte{},
		msg: msg{
			MessageType: message.TypeConnect,
			UserID:      userID,
		},
	}
}

func (c ClientConnectMessage) Bytes() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(fmt.Errorf("failed to serialize connection request: %s", err.Error()))
	}
	return b
}

func (c ClientConnectMessage) User() string {
	return c.UserID
}
