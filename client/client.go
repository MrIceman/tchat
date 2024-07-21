package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"sync"
	"tchat/internal/message"
	"tchat/internal/protocol"
	types2 "tchat/internal/types"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
	renderer *viewController
	conn     net.Conn
	id       string
	app      *app

	sendMessageChan chan []byte
	renderTextChan  chan []string
	exitCh          chan struct{}
}

func New(conn net.Conn) *Client {
	clientID := uuid.NewString()
	sendMessageCh := make(chan []byte)
	channelsJoinedCh := make(chan types2.Channel)
	renderTextCh := make(chan []string)
	exitChannelCh := make(chan struct{})
	v := newView(sendMessageCh, renderTextCh, channelsJoinedCh, exitChannelCh)
	exitCh := make(chan struct{})
	v.setUp()
	renderer := newViewController(renderTextCh, channelsJoinedCh)

	return &Client{
		conn:            conn,
		renderer:        renderer,
		id:              clientID,
		renderTextChan:  renderTextCh,
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

		c.renderTextChan <- []string{fmt.Sprintf("Connected to server as %s", connectRes.Message.(string))}
	default:
		log.Fatalf("unexpected response from server: %s", string(b))
		return
	}
}

func (c *Client) Run() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := c.app.Run(); err != nil {
			log.Fatalf("could not run application: %s", err.Error())
		}
	}()

	go func() {
		for {
			b := message.Receive(c.conn)
			c.renderer.onNewMessage(b)
		}
	}()

	go func() {
		for {
			select {
			case _ = <-c.exitCh:
				c.renderTextChan <- []string{"Disconnected from server"}
				os.Exit(0)
			case msg := <-c.sendMessageChan:
				m, err := ParseFromInput(c.id, string(msg))
				// check if m is instance of DisconnectMessage
				if _, ok := m.(protocol.DisconnectMessage); ok {
					c.exitCh <- struct{}{}
					break
				}

				if err != nil {
					c.renderTextChan <- []string{fmt.Sprintf("could not parse message: %s", err.Error())}
					c.app.lobbyView.ScrollToEnd()
				} else {
					message.Transmit(c.conn, m.Bytes())
				}
			}
		}
	}()
	wg.Wait()

}
