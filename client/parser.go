package client

import (
	"errors"
	"strings"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
	"tchat/internal/validation"
	"time"
)

const (
	channelNewUserMessagePrefix = "#newmessage#"
)

func ParseFromInput(userID, input string) (protocol.SerializableMessage, error) {
	input = strings.ReplaceAll(input, "\n", "")
	parsedInput := strings.Split(input, " ")
	msgType := parsedInput[0]
	switch msgType {
	case "/exit":
		return protocol.DisconnectMessage{}, nil
	case "/channel":
		if len(parsedInput) < 2 {
			return nil, errors.New("insufficient arguments")
		}
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
		// dont like encoding it in the message, but for now it works
		if strings.HasPrefix(input, channelNewUserMessagePrefix) {
			msgWithoutPrefix := strings.TrimPrefix(input, channelNewUserMessagePrefix)
			channelAndMsg := strings.Split(msgWithoutPrefix, "#")
			types.Message{
				UserID:      userID,
				Channel:     channelAndMsg[0],
				DisplayName: userID,
				Content:     channelAndMsg[1],
				CreatedAt:   time.Now(),
			}.MustJSON()

			return protocol.NewChannelsMessage(userID, message.TypeChannelNewMessage, []byte(input)), nil
		}
		return nil, validation.ErrMessageTypeNotImplemented
	}
}
