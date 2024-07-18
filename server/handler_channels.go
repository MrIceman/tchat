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
	default:
		log.Fatalf("unexpected message: %s", msgType)
	}
}
