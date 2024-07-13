package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"tchat/client"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatalf("could not connect to server on port 8080: %s", err.Error())
	}
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	userID := uuid.NewString()
	connect(conn, userID)
	exit := false
	reader := bufio.NewReader(os.Stdin)
	for !exit {
		var text string
		fmt.Print(fmt.Sprintf("(%s)>: ", userID))
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %s", err.Error())
		}
		switch text {
		case "exit":
			exit = true
		default:
			msg, err := client.ParseFromInput(userID, text)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			message.Transmit(conn, msg.Bytes())
		}
	}
}

func connect(conn net.Conn, userID string) {
	message.Transmit(conn, protocol.NewClientConnectMessage(userID).Bytes())
	resp, b := message.Receive(conn)

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
