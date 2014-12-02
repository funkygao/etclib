package etcdcli

import (
	"errors"
)

var (
	ErrDialFailure = errors.New("failed to connect etcd cluster")
)
