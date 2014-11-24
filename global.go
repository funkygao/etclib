package etclib

import (
	"github.com/coreos/go-etcd/etcd"
)

var (
	client  *etcd.Client
	project string
)

var (
	rootDirs = []string{
		DIR_FAE,
		DIR_ACTOR,
		DIR_MAINTAIN}
)
