package server

import (
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
		message.Transmit(conn, protocol.NewChannelsResponse(items, message.TypeChannelsGetResponse).Bytes())
		break
	default:
		log.Fatalf("unexpected message: %s", msgType)
	}
}
