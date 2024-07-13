package protocol

import "tchat/internal/message"

type ChannelsMessage struct {
	Payload interface{} `json:"payload"`
	msg
}

func NewChannelsMessage(userID string, t message.Type) *ChannelsMessage {
	return &ChannelsMessage{
		msg: msg{
			MessageType: t,
			UserID:      userID,
		},
	}
}

func NewChannelsResponse(payload interface{}, t message.Type) *ChannelsMessage {
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
	return c.msg.ToBytes()
}
