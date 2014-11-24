package etclib

const (
	DIR_PROJECT  = "proj"
	DIR_FAE      = "fae"
	DIR_ACTOR    = "actor"
	DIR_MAINTAIN = "maintain"
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
	NODE_PING_INTERVAL = 1
)
