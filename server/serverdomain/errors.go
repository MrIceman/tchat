package serverdomain

import (
	"errors"
)

var (
	ErrUserNotChannelOwner = errors.New("user is not channel owner")
)
