package etclib

import (
	"github.com/funkygao/go-etcd/etcd"
)

var (
	client  *etcd.Client
	project string

	nodes map[string]map[string]bool = make(map[string]map[string]bool)
)
