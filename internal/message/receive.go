package message

import (
	"encoding/json"
	"fmt"
	"net"
)

func Receive(conn net.Conn) ([]byte, error) {
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("could not read response from server: %s", err.Error())
	}
	b = b[:n]

	return b, nil
}

func RawFromBytes(b []byte) (map[string]interface{}, error) {
	var resp map[string]interface{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("raw messages from bytes %s, content: %s", err.Error(), string(b))
	}
	return resp, nil
}
