package client

import (
	"errors"
	"log"
	"strings"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/validation"
)

func ParseFromInput(userID, input string) (protocol.SerializableMessage, error) {
	input = strings.ReplaceAll(input, "\n", "")
	parsedInput := strings.Split(input, " ")
	msgType := parsedInput[0]
	switch msgType {
	case "channel":
		if len(parsedInput) < 2 {
			return nil, errors.New("insufficient arguments")
		}
		log.Println(parsedInput[1])
		if parsedInput[1] == "list" {
			return protocol.NewChannelsMessage(userID, message.TypeChannelsGet, nil), nil
		}
		if parsedInput[1] == "join" {
			if len(parsedInput) < 3 {
				return nil, errors.New("no channel name was provided")
			}
			channelName := parsedInput[2]
			return protocol.NewChannelsMessage(userID, message.TypeChannelsJoin, []byte(channelName)), nil
		}

		return nil, errors.New("invalid arguments")
	default:
		return nil, validation.ErrMessageTypeNotImplemented
	}
}
