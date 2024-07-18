package client

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
	"tchat/internal/message"
	"tchat/internal/protocol"
)

type Client struct {
	renderer *renderer
	conn     net.Conn
	prompter *prompter
	id       string

	sendMessageChan    chan []byte
	receiveMessageSubs []chan []byte
}

func New(conn net.Conn) *Client {
	clientID := uuid.NewString()
	sendMessageCh := make(chan []byte)
	prompterSub := make(chan []byte)
	rendererSub := make(chan []byte)
	renderer := newRenderer(rendererSub)
	prompter := NewPrompter(clientID, sendMessageCh, prompterSub)
	receiveMessageSubs := []chan []byte{rendererSub, prompterSub}

	return &Client{
		conn:               conn,
		renderer:           renderer,
		id:                 clientID,
		prompter:           prompter,
		sendMessageChan:    sendMessageCh,
		receiveMessageSubs: receiveMessageSubs,
	}
}

func (c *Client) Connect() {
	message.Transmit(c.conn, protocol.NewClientConnectMessage(c.id).Bytes())
	b := message.Receive(c.conn)
	resp := message.RawFromBytes(b)
	switch resp["type"] {
	case string(message.TypeConnectRes):
		var connectRes protocol.ServerSystemMessage
		if err := json.Unmarshal(b, &connectRes); err != nil {
			log.Fatalf("could not unmarshal connect response: %s", err.Error())
		}
		log.Printf(connectRes.Message.(string))
	default:
		log.Fatalf("unexpected response from server: %s", string(b))
		return
	}
}

func (c *Client) Run() {
	// TODO add a way to gracefully shutdown the client
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			b := message.Receive(c.conn)
			c.broadcastMessage(b)
		}
	}()

	go func() {
		c.renderer.renderMessage()
	}()

	go func() {
		c.prompter.Prompt()
	}()

	go func() {
		for {
			select {
			case msg := <-c.sendMessageChan:
				message.Transmit(c.conn, msg)
			}
		}
	}()
	log.Println("client is running")
	wg.Wait()
	log.Println("client is running")
}

func (c *Client) broadcastMessage(b []byte) {
	for _, sub := range c.receiveMessageSubs {
		sub <- b
	}
}
