package client

import (
	"encoding/json"
	"log"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
)

type renderer struct {
}

func newRenderer() *renderer {
	return &renderer{}
}

func (r *renderer) renderMessage(b []byte, msg map[string]interface{}) {
	msgType := message.Type(msg["type"].(string))
	if !msgType.IsValid() {
		log.Printf("invalid message type received: %s", msgType)
	}
	if msgType.IsChannelMsg() {
		r.renderChannelMessage(msgType, b)
	}
}

func (r *renderer) renderChannelMessage(msgType message.Type, b []byte) {
	switch msgType {
	case message.TypeChannelsGetResponse:
		c := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &c); err != nil {
			log.Fatalf("could not unmarshal response: %s", err.Error())
		}
		var channels []types.Channel
		_ = json.Unmarshal(c.Payload, &channels)
		log.Println("#### Channels ####")
		log.Printf("- %d channels found -", len(channels))
		for _, ch := range channels {
			log.Printf("\t* %s", ch.Name)
		}
	default:
		log.Printf("unhandled message type: %s", msgType)
	}
}
