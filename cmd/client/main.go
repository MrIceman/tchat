package main

import (
	"log"
	"net"
	"tchat/client"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatalf("could not connect to server on port 8080: %s", err.Error())
	}
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	c := client.New(conn)
	c.Connect()
	c.Run()
}
