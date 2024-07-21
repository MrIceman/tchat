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

	TypeChannelsGet          Type = "channel_get"
	TypeChannelsGetResponse       = "channel_get_response"
	TypeChannelsJoin         Type = "channel_join"
	TypeChannelsJoinResponse Type = "channel_join_response"
	TypeChannelNewMessage    Type = "channel_new_message"
)

var (
	allMessagesTypes = []Type{
		TypeConnect,
		TypeConnectRes,
		TypeChannelsGet,
		TypeChannelsJoin,
		TypeChannelsJoinResponse,
		TypeChannelsGetResponse,
		TypeChannelNewMessage,
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
