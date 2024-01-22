package entity

import (
	"errors"
)

var (
	ErrNoRows    = errors.New("row not found")
	ErrIDToSmall = errors.New("id too small")
)
