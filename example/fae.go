package main

import (
	"github.com/funkygao/etclib"
	log "github.com/funkygao/log4go"
)

func main() {
	etclib.Init([]string{"http://127.0.0.1:4001"}, "dw")
	log.Debug("init done")

	etclib.BootFae("localhost:9001")
	etclib.ShutdownActor("localhost:9002")

	log.Debug("getting cluster nodes...")
	nodes, err := etclib.ClusterNodes(etclib.NODE_FAE)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("fae nodes: %+v", nodes)
	<-make(chan interface{})
}
