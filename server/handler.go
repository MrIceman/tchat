package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/server/serverdomain"
)

type handler struct {
	svc   *serverdomain.UserService
	chSvc *serverdomain.ChannelService
	conns []net.Conn
	mutex sync.Mutex
}

func newHandler(chSvc *serverdomain.ChannelService, svc *serverdomain.UserService) *handler {
	return &handler{
		svc:   svc,
		chSvc: chSvc,
		mutex: sync.Mutex{},
		conns: []net.Conn{},
	}
}

func (h *handler) handleConnection(conn net.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.conns = append(h.conns, conn)
	go h.setUpConnListener(conn)
}

func (h *handler) setUpConnListener(conn net.Conn) {
	for {
		log.Println("got a message")
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Printf("err while parsing message: %s", err.Error())
			return
		}

		msgB := b[:n]
		var resp map[string]interface{}
		if err := json.Unmarshal(msgB, &resp); err != nil {
			log.Fatalf("could not unmarashal request: %s", err.Error())
		}
		msgType := message.Type(resp["type"].(string))
		if !msgType.IsValid() {
			log.Printf("invalid message type received: %s", msgType)
		}
		log.Printf("received message: %s", string(msgB))
		switch msgType {
		case message.TypeConnect:
			var connectMsg protocol.ClientConnectMessage
			if err := json.Unmarshal(msgB, &connectMsg); err != nil {
				log.Fatalf("could not unmarshal connect message: %s", err.Error())
			}
			userID := connectMsg.UserID
			if err := h.svc.SignInUser(userID); err != nil {
				log.Fatalf("could not sign user: %s", err.Error())
			}
			message.Transmit(conn, protocol.NewServerSystemMessage(fmt.Sprintf("Hello %s,\n\n%s", userID, welcomeText)).Bytes())
			log.Printf("%s connected", userID)
			return
		case message.TypeChannelsGet:
			log.Println("fetching all channels")
			items, err := h.chSvc.GetAll()
			if err != nil {
				log.Fatalf("failed to get all channels: %s", err.Error())
			}
			message.Transmit(conn, protocol.NewChannelsResponse(items, message.TypeChannelsGetResponse).Bytes())
			return
		default:
			log.Fatalf("unexpected message: %s", string(msgB))
		}
	}
}
