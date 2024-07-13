package message

import "slices"

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
