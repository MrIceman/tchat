package message

import (
	"log"
	"net"
	"time"
)

func Transmit(conn net.Conn, b []byte) {
	if _, err := conn.Write(b); err != nil {
		log.Printf("could not transmit message: %s - %d", err.Error(), time.Now().UnixMicro())
	}
}
