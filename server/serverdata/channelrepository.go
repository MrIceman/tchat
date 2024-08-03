package serverdata

import (
	"errors"
	"fmt"
	"log"
	"net"
	"slices"
	"sync"
	message2 "tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
	"time"
)

type userIP = string

type ChannelRepository struct {
	mutex                 sync.Mutex
	channelList           []*types.Channel
	channelConns          map[string][]net.Conn
	connCurrentChannelMap map[userIP]string
	connUserIDMap         map[userIP]string
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{
		mutex:        sync.Mutex{},
		channelConns: make(map[string][]net.Conn),
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
		connCurrentChannelMap: make(map[string]string),
		connUserIDMap:         make(map[string]string),
	}
}

func (cr *ChannelRepository) GetByName(channelname string) (types.Channel, error) {
	for i := range cr.channelList {
		c := cr.channelList[i]
		if c.Name == channelname {
			return *c, nil
		}
	}
	return types.Channel{}, fmt.Errorf("channel %s does not exist", channelname)
}

func (cr *ChannelRepository) GetAll() []types.Channel {
	chL := make([]types.Channel, len(cr.channelList))
	for i, c := range cr.channelList {
		chL[i] = *c
	}

	return chL
}

func (cr *ChannelRepository) OnConnectionDisconnected(conn net.Conn) error {
	usrChannel := cr.connCurrentChannelMap[conn.RemoteAddr().String()]
	usr := cr.connUserIDMap[conn.RemoteAddr().String()]
	if usrChannel == "" {
		return errors.New("no user was stored for the connection")
	}
	usrChannelName := cr.connCurrentChannelMap[conn.RemoteAddr().String()]
	if usrChannelName == "" {
		log.Printf("user was not in any usrChannelName")
		return nil
	}
	cr.mutex.Lock()
	idx := slices.IndexFunc(cr.channelList, func(channel *types.Channel) bool {
		return channel.Name == usrChannelName
	})
	if idx == -1 {
		cr.mutex.Unlock()
		return fmt.Errorf("could not find channel in channelList with name %s", usrChannelName)
	}
	channel := cr.channelList[idx]
	channel.CurrentUsers -= 1
	delete(cr.connUserIDMap, conn.RemoteAddr().String())
	delete(cr.connCurrentChannelMap, conn.RemoteAddr().String())
	allChannelCons := cr.channelConns[usrChannelName]
	connIdx := slices.IndexFunc(allChannelCons, func(c net.Conn) bool {
		return conn == c
	})

	newChannelConns := append(allChannelCons[:connIdx], allChannelCons[connIdx+1:]...)
	cr.channelConns[usrChannelName] = newChannelConns
	cr.mutex.Unlock()

	cr.sendMessageAndHandleZombieConns(protocol.NewChannelsMessage(
		usr,
		message2.TypeChannelUserDisconnectedMessage,
		[]byte(usr)).Bytes(),
		usrChannel)

	return nil
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

func (cr *ChannelRepository) OnNewUser(channelName, userID string, conn net.Conn) (*types.Channel, error) {

	for i, ch := range cr.channelList {
		if ch.Name == channelName {
			cr.mutex.Lock()
			cr.channelList[i].CurrentUsers++
			cr.channelConns[channelName] = append(cr.channelConns[channelName], conn)
			cr.connUserIDMap[conn.RemoteAddr().String()] = userID
			cr.connCurrentChannelMap[conn.RemoteAddr().String()] = channelName
			cr.mutex.Unlock()
			return ch, nil
		}
	}

	return nil, fmt.Errorf("channel with name %s not found", channelName)
}

func (cr *ChannelRepository) NewMessage(channelName string, msg types.Message) error {
	channelFound := false
	for _, ch := range cr.channelList {
		if ch.Name == channelName {
			cr.mutex.Lock()
			ch.TotalMessages++
			cr.mutex.Unlock()
			channelFound = true
		}
	}

	if !channelFound {
		return fmt.Errorf("channel with name %s not found", channelName)
	}

	cr.sendMessageAndHandleZombieConns(protocol.NewChannelsMessage(msg.UserID,
		message2.TypeChannelNewMessage, msg.MustJSON()).Bytes(), channelName)

	return nil
}

func (cr *ChannelRepository) sendMessageAndHandleZombieConns(b []byte, channelName string) {
	channelConns := cr.channelConns[channelName]
	zombieConns := message2.Broadcast(channelConns, b)

	if len(zombieConns) > 0 {
		cr.mutex.Lock()
		for _, conn := range zombieConns {
			log.Printf("conn %s not reachable, removing from channel %s", conn.RemoteAddr(), channelName)
			for i, c := range channelConns {
				if c == conn {
					channelConns = append(channelConns[:i], channelConns[i+1:]...)
				}
			}
		}
		// probably not necessary
		cr.channelConns[channelName] = channelConns
		cr.mutex.Unlock()
	}
}

func (cr *ChannelRepository) Delete(channelname string) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	for i := range cr.channelList {
		c := cr.channelList[i]
		if c.Name == channelname {
			cr.sendMessageAndHandleZombieConns(protocol.NewChannelsMessage(
				"system",
				message2.TypeChannelMustLeave, nil).Bytes(),
				channelname)
			cr.channelList = append(cr.channelList[:i], cr.channelList[i+1:]...)
			return nil
		}
	}

	return errors.New("channel could not be found")
}
