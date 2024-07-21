package client

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"io"
	"log"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
)

type renderer struct {
	receiveMessageChan chan []byte
	writer             io.Writer
}

func newRenderer(receiveMessageChan chan []byte, writer io.Writer) *renderer {
	return &renderer{
		receiveMessageChan: receiveMessageChan,
		writer:             writer,
	}
}

func (r *renderer) renderMessage(b []byte) {
	msg := message.RawFromBytes(b)
	msgType := message.Type(msg["type"].(string))
	if !msgType.IsValid() {
		fmt.Fprintln(r.writer, "invalid message type received: %s", msgType)
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
		fmt.Fprintln(r.writer, "#### Channels ####")
		fmt.Fprintln(r.writer, fmt.Sprintf("- %d channels found -", len(channels)))
		for _, ch := range channels {
			fmt.Fprintln(r.writer, fmt.Sprintf("\t* %s (%d online)", ch.Name, ch.CurrentUsers))
		}
	case message.TypeChannelsJoinResponse:
		c := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &c); err != nil {
			log.Fatalf("could not unmarshal response: %s", err.Error())
		}
		channel := types.Channel{}
		_ = json.Unmarshal(c.Payload, &channel)
		fmt.Fprintln(r.writer, fmt.Sprintf("#### Joined Channel %s - There are currently %d users online ####", channel.Name, channel.CurrentUsers))
		fmt.Fprintln(r.writer, channel.WelcomeMessage)

	default:
		fmt.Fprintf(r.writer, "unhandled message type: %s", msgType)
	}

	r.writer.(*tview.TextView).ScrollToEnd()
}
