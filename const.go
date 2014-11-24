package etclib

const (
	DIR_PROJECT = "proj"
	DIR_NODES   = "node"

	DIR_MAINTAIN = "down"
)

const (
	NODE_FAE   = "fae"
	NODE_ACTOR = "actor"
)

const (
	OP_BOOT = iota + 1
	OP_SHUTDOWN
)

const (
	STATUS_ALIVE = "alive"
	STATUS_DIED  = "died"
)

const (
	NODE_PING_INTERVAL = 60 // in seconds
)
