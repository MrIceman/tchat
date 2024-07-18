package serverdata

import (
	"fmt"
	"slices"
	"sync"
	"tchat/internal/types"
	"time"
)

type ChannelRepository struct {
	mutex       sync.Mutex
	channelList []types.Channel
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{
		mutex: sync.Mutex{},
		channelList: []types.Channel{
			{
				Name:          "general",
				Owner:         "system",
				CreatedAt:     time.Now(),
				CurrentUsers:  0,
				TotalMessages: 0,
			},
		},
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
