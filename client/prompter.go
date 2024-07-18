package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type prompter struct {
	clientID      string
	currentPrefix string
	sendMessageCh chan []byte
	readMessageCh chan []byte
}

func NewPrompter(clientID string, sendMessageCh chan []byte, readMessageCh chan []byte) *prompter {
	return &prompter{
		clientID:      clientID,
		currentPrefix: fmt.Sprintf("(%s)>: ", clientID),
		sendMessageCh: sendMessageCh,
		readMessageCh: readMessageCh,
	}
}

func (p *prompter) startListening() {
	for {
		select {
		case msg := <-p.readMessageCh:
			_ = msg
		}
	}
}

func (p *prompter) Prompt() {
	exit := false
	reader := bufio.NewReader(os.Stdin)
	go p.startListening()

	for !exit {
		var text string
		fmt.Print(p.currentPrefix)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %s", err.Error())
		}
		switch text {
		case "exit":
			exit = true
			continue
		default:
			msg, err := ParseFromInput(p.clientID, text)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			p.sendMessageCh <- msg.Bytes()
			continue
		}
	}
}
