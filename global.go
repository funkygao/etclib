package etclib

import (
	"errors"
)

var (
	store Backend

	nodes map[string]map[string]bool = make(map[string]map[string]bool)
)

var (
	ErrInvalidService = errors.New("Invalid service type")
)
