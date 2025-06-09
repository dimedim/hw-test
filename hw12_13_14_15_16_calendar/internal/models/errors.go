package models

import "errors"

var (
	ErrEventNotExists = errors.New("event not exists")
	ErrEventNotSet    = errors.New("event not set")
)
