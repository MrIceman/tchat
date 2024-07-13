package types

import "time"

type User struct {
	UserID     string    `json:"UserID"`
	LoggedInAt time.Time `json:"CreatedAt"`
}
