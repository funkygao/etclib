package etclib

import (
	"time"
)

type Client interface {
	DialTimeout(servers []string, timeout time.Duration) error
	PrepareDirs(nodeDirs ...string) error
	Close()

	BootFae(addr string)
	ShutdownFae(addr string)
	BootActor(addr string)
	ShutdownActor(addr string)
	WatchFaeNodes() (ch chan NodeEvent)
	ClusterNodes(nodeType string) ([]string, error)

	MaintainInfo() ([]string, error)
	WatchMaintain() (ch chan MaintainEvent)
}
