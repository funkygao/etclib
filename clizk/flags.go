package clizk

import (
	"github.com/samuel/go-zookeeper/zk"
)

const (
	FlagNormal          = 0 // persistent
	FlagService         = zk.FlagEphemeral
	FlagSequence        = zk.FlagSequence
	FlagServiceSequence = zk.FlagEphemeral | zk.FlagSequence
)
