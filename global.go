package etclib

import (
	"github.com/coreos/go-etcd/etcd"
)

var (
	client  *etcd.Client
	project string
)
