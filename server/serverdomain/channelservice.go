package serverdomain

import (
	"net"
	"tchat/internal/types"
	"tchat/server/serverdata"
	"time"
)

type ChannelService struct {
	userRepository *serverdata.UserRepository
	repository     *serverdata.ChannelRepository
}

func NewChannelService(userRepository *serverdata.UserRepository, channelRepository *serverdata.ChannelRepository) *ChannelService {
	return &ChannelService{
		userRepository: userRepository,
		repository:     channelRepository,
	}
}

func (cs *ChannelService) CreateChannel(userID, channelName string) error {
	c := types.Channel{
		Name:      channelName,
		Owner:     userID,
		CreatedAt: time.Now(),
	}

	return cs.repository.CreateChannel(c)
}

func (cs *ChannelService) GetAll() ([]types.Channel, error) {
	return cs.repository.GetAll(), nil
}

func (cs *ChannelService) JoinChannel(userID, channelName string, conn net.Conn) (*types.Channel, error) {
	return cs.repository.OnNewUser(channelName, userID, conn)
}

func (cs *ChannelService) SendToChannel(channelName string, msg types.Message) error {
	return cs.repository.NewMessage(channelName, msg)
}

func (cs *ChannelService) UserDisconnected(conn net.Conn) error {
	return cs.repository.OnConnectionDisconnected(conn)
}

func (cs *ChannelService) DeleteChannel(id string, channelname string) error {
	channel, err := cs.repository.GetByName(channelname)
	if err != nil {
		return err
	}

	if channel.Owner != id {
		return ErrUserNotChannelOwner
	}
	return cs.repository.Delete(channelname)
}
