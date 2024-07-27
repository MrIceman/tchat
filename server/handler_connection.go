package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"tchat/internal/message"
	"tchat/internal/protocol"
)

func (h *handler) handleConnectionMessage(conn net.Conn, msgType message.Type, b []byte) {
	switch msgType {
	case message.TypeConnect:
		var connectMsg protocol.ClientConnectMessage
		if err := json.Unmarshal(b, &connectMsg); err != nil {
			log.Fatalf("could not unmarshal connect message: %s", err.Error())
		}
		userID := connectMsg.UserID
		if err := h.userSvc.SignInUser(userID); err != nil {
			log.Fatalf("could not sign user: %s", err.Error())
		}
		message.Transmit(conn, protocol.NewServerSystemMessage(fmt.Sprintf("Hello %s,\n\n%s", userID, welcomeText)).Bytes())
		log.Printf("%s (%s) connected", conn.RemoteAddr(), userID)
		return
	default:
		log.Fatalf("unexpected message: %s", string(msgType))
	}
}
