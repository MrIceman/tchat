package main

import (
	"log"
	"net"
	"tchat/client"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatalf("could not connect to server on port 8080: %s", err.Error())
	}
	c := client.New(conn)
	c.Connect()
	c.Run()
}
