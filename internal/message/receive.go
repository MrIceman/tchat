package message

import (
	"encoding/json"
	"log"
	"net"
)

func Receive(conn net.Conn) []byte {
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		log.Fatalf("could not read response from server: %s", err.Error())
	}
	b = b[:n]

	return b
}

func RawFromBytes(b []byte) map[string]interface{} {
	var resp map[string]interface{}
	if err := json.Unmarshal(b, &resp); err != nil {
		log.Fatalf("could not unmarshal response: %s", err.Error())
	}
	return resp
}
