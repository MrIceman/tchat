package validation

import "errors"

var (
	ErrInvalidMessageType        = errors.New("invalid message type")
	ErrMessageTypeNotImplemented = errors.New("message type not implemented")
)
