package common

import "errors"

var (
	ErrNoItem = errors.New("items must have at least one item")
)
