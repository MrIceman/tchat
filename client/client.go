package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
	"tchat/internal/message"
	"tchat/internal/protocol"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
	renderer *renderer
	conn     net.Conn
	id       string
	app      *app

	sendMessageChan    chan []byte
	receiveMessageSubs []chan []byte
	exitCh             chan struct{}
}

func New(conn net.Conn) *Client {
	clientID := uuid.NewString()
	sendMessageCh := make(chan []byte)
	rendererSub := make(chan []byte)
	v := newView(sendMessageCh)
	exitCh := make(chan struct{})
	v.setUp()
	renderer := newRenderer(rendererSub, v.UI())

	return &Client{
		conn:            conn,
		renderer:        renderer,
		id:              clientID,
		sendMessageChan: sendMessageCh,
		app:             v,
		exitCh:          exitCh,
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

		fmt.Fprintln(c.renderer.writer, fmt.Sprintf("Connected to server as %s", connectRes.Message.(string)))
	default:
		log.Fatalf("unexpected response from server: %s", string(b))
		return
	}
}

func (c *Client) Run() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			<-c.exitCh
			c.app.application.Stop()
		}
	}()
	go func() {
		if err := c.app.Run(); err != nil {
			log.Fatalf("could not run application: %s", err.Error())
		}
	}()

	go func() {
		for {
			b := message.Receive(c.conn)
			c.renderer.renderMessage(b)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-c.sendMessageChan:
				m, err := ParseFromInput(c.id, string(msg))
				// check if m is instance of DisconnectMessage
				if _, ok := m.(protocol.DisconnectMessage); ok {
					c.exitCh <- struct{}{}
					break
				}

				if err != nil {
					fmt.Fprintln(c.app.UI(), fmt.Sprintf("could not parse message: %s", err.Error()))
					c.app.textView.ScrollToEnd()
				} else {
					message.Transmit(c.conn, m.Bytes())
				}
			}
		}
	}()
	wg.Wait()

}

func (c *Client) broadcastMessage(b []byte) {
	for _, sub := range c.receiveMessageSubs {
		sub <- b
	}
}
