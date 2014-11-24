package etclib

const (
	DIR_PROJECT = "proj"
	DIR_NODES   = "node"

	DIR_MAINTAIN = "maintain"
)

const (
	NODE_FAE   = "fae"
	NODE_ACTOR = "act"
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
	NODE_PING_INTERVAL = 600 // in seconds
)
