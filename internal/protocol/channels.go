package protocol

import (
	"encoding/json"
	"tchat/internal/message"
)

type ChannelsMessage struct {
	Payload []byte `json:"payload"`
	msg
}

func NewChannelsMessage(userID string, t message.Type, payload []byte) *ChannelsMessage {
	return &ChannelsMessage{
		Payload: payload,
		msg: msg{
			MessageType: t,
			UserID:      userID,
		},
	}
}

func NewChannelsResponse(payload []byte, t message.Type) *ChannelsMessage {
	return &ChannelsMessage{
		Payload: payload,
		msg: msg{
			MessageType: t,
		},
	}
}

func (c ChannelsMessage) User() string {
	return c.UserID
}

func (c ChannelsMessage) Bytes() []byte {
	b, _ := json.Marshal(c)
	return b
}
