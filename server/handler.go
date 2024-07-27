package server

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"tchat/internal/message"
	"tchat/server/serverdomain"
	"time"
)

type handler struct {
	userSvc *serverdomain.UserService
	chSvc   *serverdomain.ChannelService
	conns   []net.Conn
	mutex   sync.Mutex
}

func newHandler(chSvc *serverdomain.ChannelService, svc *serverdomain.UserService) *handler {
	return &handler{
		userSvc: svc,
		chSvc:   chSvc,
		mutex:   sync.Mutex{},
		conns:   []net.Conn{},
	}
}

func (h *handler) handleConnection(conn net.Conn) {
	log.Printf("handling connection")
	h.mutex.Lock()
	h.conns = append(h.conns, conn)
	h.mutex.Unlock()
	go h.setUpConnListener(conn)
}

func (h *handler) setUpConnListener(conn net.Conn) {
	defer conn.Close()
	for {
		log.Println("waiting for message from client...")
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection %s closed by client. %d", conn.RemoteAddr(), time.Now().UnixMicro())
				return
			}

			log.Printf("err while parsing message: %s", err.Error())
			continue
		}

		b = b[:n]
		go func(b []byte) {
			if os.Getenv("debug") == "true" {
				log.Printf("received message: %s", string(b))
			}
			var msg map[string]interface{}
			if err := json.Unmarshal(b, &msg); err != nil {
				log.Fatalf("could not unmarshal request: %s", err.Error())
			}
			msgType := message.Type(msg["type"].(string))
			if !msgType.IsValid() {
				log.Printf("invalid message type received: %s", msgType)
			}

			if msgType.IsChannelMsg() {
				h.handleChannelMessage(conn, msgType, b)
			}

			if msgType.IsConnectMsg() {
				h.handleConnectionMessage(conn, msgType, b)
			}
		}(b)
	}
}
