package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	UserID      string    `json:"userID"`
	DisplayName string    `json:"displayName"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (m Message) MustJSON() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("failed to serialize message: %s", err.Error()))
	}
	return b
}
