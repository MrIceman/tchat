package protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"tchat/internal/message"
)

type SerializableMessage interface {
	Bytes() []byte

	User() string
}

type msg struct {
	UserID      string       `json:"user_id"`
	MessageType message.Type `json:"type"`
}

func (c msg) ToBytes() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(fmt.Errorf("failed to serialize connection request: %s", err.Error()))
	}
	log.Printf("serialized message: %s", string(b))
	return b
}
