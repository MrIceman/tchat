package client

import (
	"encoding/json"
	"fmt"
	"log"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
)

type viewController struct {
	onChannelJoinCh chan types.Channel
	renderTextCh    chan []string
}

func newViewController(renderTextCh chan []string, onChannelJoinCh chan types.Channel) *viewController {
	return &viewController{
		renderTextCh:    renderTextCh,
		onChannelJoinCh: onChannelJoinCh,
	}
}

func (r *viewController) onNewMessage(b []byte) {
	msg := message.RawFromBytes(b)
	msgType := message.Type(msg["type"].(string))
	if !msgType.IsValid() {
		r.renderTextCh <- []string{fmt.Sprintf("invalid message type received: %s", msgType)}
	}
	if msgType.IsChannelMsg() {
		r.renderChannelMessage(msgType, b)
	}
}

func (r *viewController) renderChannelMessage(msgType message.Type, b []byte) {
	switch msgType {
	case message.TypeChannelsGetResponse:
		c := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &c); err != nil {
			log.Fatalf("could not unmarshal response: %s", err.Error())
		}
		var channels []types.Channel
		_ = json.Unmarshal(c.Payload, &channels)
		r.renderTextCh <- []string{"#### Channels ####", fmt.Sprintf("- %d channels found -", len(channels))}
		for _, ch := range channels {
			r.renderTextCh <- []string{fmt.Sprintf("\t* %s (%d online)", ch.Name, ch.CurrentUsers)}
		}
	case message.TypeChannelsJoinResponse:
		c := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &c); err != nil {
			log.Fatalf("could not unmarshal response: %s", err.Error())
		}
		channel := types.Channel{}
		_ = json.Unmarshal(c.Payload, &channel)
		r.onChannelJoinCh <- channel
		r.renderTextCh <- []string{fmt.Sprintf("#### Joined Channel %s - There are currently %d users online ####", channel.Name, channel.CurrentUsers)}
		r.renderTextCh <- []string{"#### Type /leave to leave the channel ####"}
		r.renderTextCh <- []string{"#### Type /channels to see all available channels ####"}
		r.renderTextCh <- []string{"#### Type /users to see all users in the channel ####"}
		r.renderTextCh <- []string{"#### Type /msg <user> <message> to send a private message ####"}
		r.renderTextCh <- []string{"----------------------------------------------------------", fmt.Sprintf("Channel Message: %s", channel.WelcomeMessage)}
	default:
		log.Fatalf("unexpected message type: %s", msgType)
	}
}
