package serverdata

import (
	"fmt"
	"net"
	"slices"
	"sync"
	message2 "tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
	"time"
)

type ChannelRepository struct {
	mutex        sync.Mutex
	channelList  []*types.Channel
	channelConns map[string][]net.Conn
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

func (cr *ChannelRepository) OnNewUser(channelName string, conn net.Conn) (*types.Channel, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	for i, ch := range cr.channelList {
		if ch.Name == channelName {
			cr.channelList[i].CurrentUsers++
			return ch, nil
		}
	}
	cr.channelConns[channelName] = append(cr.channelConns[channelName], conn)

	return nil, fmt.Errorf("channel with name %s not found", channelName)
}

func (cr *ChannelRepository) NewMessage(channelName string, msg types.Message) error {
	for _, ch := range cr.channelList {
		if ch.Name == channelName {
			cr.mutex.Lock()
			ch.TotalMessages++
			cr.mutex.Unlock()
			return nil
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(len(cr.channelConns[channelName]))
	for _, conn := range cr.channelConns[channelName] {
		go func(wg *sync.WaitGroup) {
			message2.Transmit(conn,
				protocol.NewChannelsMessage(msg.UserID, message2.TypeChannelNewMessage, msg.MustJSON()).Bytes())
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	return fmt.Errorf("channel with name %s not found", channelName)
}
