package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"tchat/internal/message"
	"tchat/internal/protocol"
)

type Client struct {
	conn net.Conn
	id   string
}

func New(conn net.Conn) *Client {
	return &Client{
		conn: conn,
		id:   uuid.NewString(),
	}
}

func (c *Client) Connect() {
	message.Transmit(c.conn, protocol.NewClientConnectMessage(c.id).Bytes())
	resp, b := message.Receive(c.conn)

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
	exit := false
	reader := bufio.NewReader(os.Stdin)
	for !exit {
		var text string
		fmt.Print(fmt.Sprintf("(%s)>: ", c.id))
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %s", err.Error())
		}
		switch text {
		case "exit":
			exit = true
		default:
			msg, err := ParseFromInput(c.id, text)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			message.Transmit(c.conn, msg.Bytes())
		}
	}
}
