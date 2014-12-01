package etclib

var (
	client  Client
	project string

	nodes map[string]map[string]bool = make(map[string]map[string]bool)
)
