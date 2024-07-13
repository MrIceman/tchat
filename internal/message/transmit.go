package message

import (
	"log"
	"net"
)

func Transmit(conn net.Conn, b []byte) {
	if _, err := conn.Write(b); err != nil {
		log.Fatalf("could not transmit message: %s", err.Error())
	}
}
