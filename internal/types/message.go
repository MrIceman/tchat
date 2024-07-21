package types

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	Content     string `json:"content"`
	CreatedAt   string `json:"createdAt"`
}

func (m Message) MustJSON() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("failed to serialize message: %s", err.Error()))
	}
	return b
}
