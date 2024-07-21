package message

import (
	"log"
	"net"
	"os"
	"time"
)

func Transmit(conn net.Conn, b []byte) {
	if os.Getenv("debug") == "true" {
		log.Printf("transmitting message: %s - %d", string(b), time.Now().UnixMicro())
	}
	if _, err := conn.Write(b); err != nil {
		log.Printf("could not transmit message: %s - %d", err.Error(), time.Now().UnixMicro())
	}
}
