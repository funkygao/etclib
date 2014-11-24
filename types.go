package etclib

import (
	"time"
)

type Node struct {
	addr         string
	lastSyncTime time.Time
}
