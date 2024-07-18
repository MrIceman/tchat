package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type prompter struct {
	clientID      string
	currentPrefix string
	sendMessageCh chan []byte
	readMessageCh chan []byte
	terminalMutex *sync.Mutex
}

func newPrompter(terminalMutex *sync.Mutex, clientID string, sendMessageCh chan []byte, readMessageCh chan []byte) *prompter {
	return &prompter{
		terminalMutex: terminalMutex,
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

// TODO this should be refactored as there can only be one coroutine that controls the terminal
func (p *prompter) Prompt() {
	reader := bufio.NewReader(os.Stdin)

	p.terminalMutex.Lock()
	var text string
	fmt.Print(p.currentPrefix)
	text, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalf("error reading input: %s", err.Error())
	}

	msg, err := ParseFromInput(p.clientID, text)
	if err != nil {
		p.terminalMutex.Unlock()
		log.Println(err.Error())
	} else {
		p.terminalMutex.Unlock()
		p.sendMessageCh <- msg.Bytes()
	}
}
