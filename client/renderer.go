package client

import (
	"encoding/json"
	"log"
	"sync"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
)

type renderer struct {
	receiveMessageChan chan []byte
	terminalMutex      *sync.Mutex
}

func newRenderer(terminalMutex *sync.Mutex, receiveMessageChan chan []byte) *renderer {
	return &renderer{
		receiveMessageChan: receiveMessageChan,
		terminalMutex:      terminalMutex,
	}
}

func (r *renderer) renderMessage() {
	for {
		b := <-r.receiveMessageChan
		log.Println("waiting for mutex to be unlocked")
		r.terminalMutex.Lock()
		log.Println("mutex unlocked")
		msg := message.RawFromBytes(b)
		msgType := message.Type(msg["type"].(string))
		if !msgType.IsValid() {
			log.Printf("invalid message type received: %s", msgType)
		}
		if msgType.IsChannelMsg() {
			r.renderChannelMessage(msgType, b)
		}
		r.terminalMutex.Unlock()
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
			log.Printf("\t* %s (%d online)", ch.Name, ch.CurrentUsers)
		}
	case message.TypeChannelsJoinResponse:
		c := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &c); err != nil {
			log.Fatalf("could not unmarshal response: %s", err.Error())
		}
		channel := types.Channel{}
		_ = json.Unmarshal(c.Payload, &channel)
		log.Printf("#### Joined Channel %s - There are currently %d users online ####", channel.Name, channel.CurrentUsers)
		log.Println(channel.WelcomeMessage)
	default:
		log.Printf("unhandled message type: %s", msgType)
	}
}
