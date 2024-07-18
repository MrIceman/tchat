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
	channelList []*types.Channel
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{
		mutex: sync.Mutex{},
		channelList: []*types.Channel{
			{
				Name:           "general",
				Owner:          "system",
				CreatedAt:      time.Now(),
				CurrentUsers:   0,
				TotalMessages:  0,
				WelcomeMessage: "Welcome to the Jungle",
			},
		},
	}
}

func (cr *ChannelRepository) GetAll() []types.Channel {
	chL := make([]types.Channel, len(cr.channelList))
	for i, c := range cr.channelList {
		chL[i] = *c
	}

	return chL
}

func (cr *ChannelRepository) CreateChannel(c types.Channel) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	if slices.ContainsFunc(cr.channelList, func(channel *types.Channel) bool {
		return channel.Name == c.Name
	}) {
		return fmt.Errorf("channel with name %s exists already", c.Name)
	}
	cr.channelList = append(cr.channelList, &c)

	return nil
}

func (cr *ChannelRepository) OnNewUser(channelName string) (*types.Channel, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	for i, ch := range cr.channelList {
		if ch.Name == channelName {
			cr.channelList[i].CurrentUsers++
			return ch, nil
		}
	}

	return nil, fmt.Errorf("channel with name %s not found", channelName)
}
