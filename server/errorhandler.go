package server

import (
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/server/serverdomain"
)

var errorMessageMap = map[error]protocol.ChannelsMessage{
	serverdomain.ErrUserNotChannelOwner: *protocol.NewChannelsResponse([]byte("user is not channel owner"), message.TypeChannelDeleteFailedResponse),
}

func GetMessageForError(err error) protocol.ChannelsMessage {
	msg, ok := errorMessageMap[err]
	if ok {
		return msg
	}

	return *protocol.NewChannelsResponse([]byte(err.Error()), message.TypeError)
}
