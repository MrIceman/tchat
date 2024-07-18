package message

import (
	"slices"
	"strings"
)

type Type string

const (
	TypeConnect    Type = "connect"
	TypeConnectRes Type = "connect_res"
	TypeDisconnect Type = "disconnect"

	TypeChannelsGet         Type = "get_channels"
	TypeChannelsGetResponse      = "get_channels_response"
)

var (
	allMessagesTypes = []Type{
		TypeConnect,
		TypeConnectRes,
		TypeChannelsGet,
		TypeDisconnect,
	}
)

func (t Type) IsValid() bool {

	return slices.Contains(allMessagesTypes, t)
}

func (t Type) IsChannelMsg() bool {
	return strings.Contains(string(t), "channel")
}

func (t Type) IsConnectMsg() bool {
	return strings.Contains(string(t), "connect")
}
