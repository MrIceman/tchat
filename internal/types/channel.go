package types

import "time"

type Channel struct {
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	CreatedAt     time.Time `json:"createdAt"`
	CurrentUsers  int       `json:"currentUsers"`
	TotalMessages int       `json:"totalMessages"`
	Password      *string
}
