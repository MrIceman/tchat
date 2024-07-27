package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"tchat/internal/message"
	"tchat/internal/protocol"
	"tchat/internal/types"
	"time"
)

// TODO currently we're transmitting to the client within the handler but also now within the channel repository
// since the handler will not return anything as there is no client facing API, we should probably move all the transmission
// to a single place, probably within the repository
func (h *handler) handleChannelMessage(conn net.Conn, msgType message.Type, b []byte) {
	switch msgType {
	case message.TypeChannelsGet:
		log.Println("fetching all channels")
		items, err := h.chSvc.GetAll()
		if err != nil {
			log.Fatalf("failed to get all channels: %s", err.Error())
		}
		b, _ := json.Marshal(items)
		message.Transmit(conn, protocol.NewChannelsResponse(b, message.TypeChannelsGetResponse).Bytes())
		break
	case message.TypeChannelsJoin:
		channelMsg := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &channelMsg); err != nil {
			log.Fatalf("could not unmarshal channel message: %s", err.Error())
		}
		channelName := string(channelMsg.Payload)
		ch, err := h.chSvc.JoinChannel(channelMsg.User(), channelName, conn)
		if err != nil {
			log.Printf("could not join channel: %s", err.Error())
			break
		}
		b, _ := json.Marshal(ch)
		message.Transmit(conn, protocol.NewChannelsResponse(b, message.TypeChannelsJoinResponse).Bytes())
		time.Sleep(1 * time.Second)
		if err := h.chSvc.SendToChannel(channelName, types.Message{
			UserID:      "system",
			DisplayName: "system",
			Content:     fmt.Sprintf("%s has joined the channel. Say a warm hello!", channelMsg.User()),
			CreatedAt:   time.Now(),
		}); err != nil {
			log.Printf("could not send message to channel: %s", err.Error())
		}
	case message.TypeChannelNewMessage:
		channelMsg := protocol.ChannelsMessage{}
		if err := json.Unmarshal(b, &channelMsg); err != nil {
			log.Printf("could not unmarshal channel message: %s", err.Error())
			break
		}
		msg := types.Message{}
		if err := json.Unmarshal(channelMsg.Payload, &msg); err != nil {
			log.Printf("could not unmarshal message: %s:", err.Error(), string(channelMsg.Payload))
			break
		}
		if err := h.chSvc.SendToChannel(msg.Channel, msg); err != nil {
			log.Printf("could not send message to channel: %s", err.Error())
			break
		}
	case message.TypeChannelsCreate:
		channelMsg, _ := parseChannelMessage(b)
		channel := types.Channel{}
		if err := json.Unmarshal(channelMsg.Payload, &channel); err != nil {
			log.Printf("could not unmarshal channel: %s", err.Error())
			break
		}
		if err := h.chSvc.CreateChannel(channelMsg.User(), channel.Name); err != nil {
			log.Printf("could not create channel: %s", err.Error())
			break
		}
		if err := message.
			Transmit(conn, protocol.NewChannelsResponse([]byte{}, message.TypeChannelsCreateResponse).Bytes()); err != nil {
			log.Printf("could not transmit message: %s", err.Error())
		}

		break
	case message.TypeChannelLeave:
		_, _ = parseChannelMessage(b)
		_ = message.Transmit(conn, protocol.NewChannelsResponse(nil, message.TypeChannelsLeaveResponse).Bytes())

		if err := h.chSvc.UserDisconnected(conn); err != nil {
			log.Printf("could not transmit message: %s", err.Error())
			break
		}

	default:
		log.Fatalf("unexpected message: %s", msgType)
	}
}

func parseChannelMessage(b []byte) (*protocol.ChannelsMessage, error) {
	channelMsg := protocol.ChannelsMessage{}
	if err := json.Unmarshal(b, &channelMsg); err != nil {
		log.Printf("could not unmarshal channel message: %s", err.Error())
		return nil, fmt.Errorf("could not parse channel msg: %s", err.Error())
	}
	return &channelMsg, nil
}
