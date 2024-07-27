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
	ctx      *clientContext
	renderer *viewController
	conn     net.Conn
	id       string

	app             *app
	sendMessageChan chan []byte
	renderTextChan  chan []string

	exitCh chan struct{}
}

func New(conn net.Conn) *Client {
	ctx := newClientContext()
	clientID := uuid.NewString()
	sendMessageCh := make(chan []byte)
	channelsJoinedCh := make(chan types2.Channel)
	renderTextCh := make(chan []string)
	exitChannelCh := make(chan struct{})
	exitCh := make(chan struct{})
	v := newView(ctx, sendMessageCh, renderTextCh, channelsJoinedCh, exitChannelCh)
	renderer := newViewController(ctx, renderTextCh, channelsJoinedCh, exitChannelCh)
	v.setUp()

	return &Client{
		conn:            conn,
		renderer:        renderer,
		id:              clientID,
		renderTextChan:  renderTextCh,
		sendMessageChan: sendMessageCh,
		app:             v,
		exitCh:          exitCh,
		ctx:             ctx,
	}
}

func (c *Client) Connect() {
	_ = message.Transmit(c.conn, protocol.NewClientConnectMessage(c.id).Bytes())
	b, _ := message.Receive(c.conn)
	resp, _ := message.RawFromBytes(b)
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
			b, err := message.Receive(c.conn)
			if err != nil {
				if err.Error() == "EOF" {
					c.exitCh <- struct{}{}
					break
				}
				c.renderTextChan <- []string{fmt.Sprintf("could not read message from server: %s", err.Error())}
				continue
			}
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
				m, err := ParseFromInput(c.ctx, c.id, string(msg))
				// check if m is instance of DisconnectMessage
				if _, ok := m.(protocol.DisconnectMessage); ok {
					c.exitCh <- struct{}{}
					break
				}

				if err != nil {
					c.renderTextChan <- []string{fmt.Sprintf("could not parse message: %s", err.Error())}
					c.app.lobbyView.ScrollToEnd()
				} else {
					_ = message.Transmit(c.conn, m.Bytes())
				}
			}
		}
	}()
	wg.Wait()

}
