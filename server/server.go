package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"tchat/server/serverdata"
	"tchat/server/serverdomain"
)

func Start() {
	userRepo := serverdata.NewUserRepository()
	channelRepo := serverdata.NewChannelRepository()

	h := newHandler(
		serverdomain.NewChannelService(userRepo, channelRepo),
		serverdomain.NewService(userRepo))

	runServer(h)
}

func runServer(h *handler) {
	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			go h.handleConnection(conn)
		}
	}()
	log.Println("t-chat server started on port 8080")
	wg.Wait()
}
