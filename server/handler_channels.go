package server

import (
	"encoding/json"
	"log"
	"net"
	"tchat/internal/message"
	"tchat/internal/protocol"
)

func (h *handler) handleChannelMessage(conn net.Conn, msgType message.Type, b []byte) {
	switch msgType {
	case message.TypeChannelsGet:
		log.Println("fetching all channels")
		items, err := h.chSvc.GetAll()
		if err != nil {
			log.Fatalf("failed to get all channels: %s", err.Error())
		}
		b, _ := json.Marshal(items)
		message.Transmit(conn, protocol.NewChannelsResponse(b, message.TypeChannelsGetResponse).Bytes())
		break
	case message.TypeChannelsJoin:
		channelMsg := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &channelMsg); err != nil {
			log.Fatalf("could not unmarshal channel message: %s", err.Error())
		}
		channelName := string(channelMsg.Payload)
		ch, err := h.chSvc.JoinChannel(channelMsg.User(), channelName)
		if err != nil {
			log.Fatalf("could not join channel: %s", err.Error())
		}
		b, _ := json.Marshal(ch)
		message.Transmit(conn, protocol.NewChannelsResponse(b, message.TypeChannelsJoinResponse).Bytes())
	default:
		log.Fatalf("unexpected message: %s", msgType)
	}
}
