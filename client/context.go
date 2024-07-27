package client

import (
	"errors"
	"tchat/internal/types"
)

const (
	currentChannelKey = "current_channel"
)

type clientContext struct {
	currentChannel *types.Channel
}

func newClientContext() *clientContext {
	return &clientContext{}
}

func (cc *clientContext) SetChannel(c *types.Channel) error {
	if cc.currentChannel != nil {
		return errors.New("user is already in a channel")
	}
	cc.currentChannel = c
	return nil
}

func (cc *clientContext) RemoveChannel() error {
	if cc.currentChannel == nil {
		return errors.New("user has no channel")
	}
	cc.currentChannel = nil
	return nil
}
