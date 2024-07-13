package serverdata

import (
	"fmt"
	"slices"
	"sync"
	"tchat/internal/types"
)

type ChannelRepository struct {
	mutex       sync.Mutex
	channelList []types.Channel
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{
		mutex:       sync.Mutex{},
		channelList: []types.Channel{},
	}
}

func (cr *ChannelRepository) GetAll() []types.Channel {
	return cr.channelList
}

func (cr *ChannelRepository) CreateChannel(c types.Channel) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	
	if slices.Contains(cr.channelList, c) {
		return fmt.Errorf("channel with name %s exists already", c.Name)
	}
	cr.channelList = append(cr.channelList, c)

	return nil
}
