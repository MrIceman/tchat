package message

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func Transmit(conn net.Conn, b []byte) error {
	if os.Getenv("debug") == "true" {
		log.Printf("transmitting message: %s - %d", string(b), time.Now().UnixMicro())
	}
	if _, err := conn.Write(b); err != nil {
		return fmt.Errorf("could not transmit message: %s - %d", err.Error(), time.Now().UnixMicro())
	}
	return nil
}
