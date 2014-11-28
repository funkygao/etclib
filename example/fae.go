package main

import (
	"fmt"
	"github.com/funkygao/etclib"
)

func main() {
	if err := etclib.Dial([]string{"http://127.0.0.1:4001"}, "dw"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("init done")

	etclib.BootFae("localhost:9001")
	etclib.ShutdownActor("localhost:9002")

	fmt.Println("getting cluster nodes...")
	nodes, err := etclib.ClusterNodes(etclib.NODE_FAE)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("fae nodes: %+v\n", nodes)
	fmt.Println("watch for nodes changes...")

	for evt := range etclib.WatchFaeNodes() {
		fmt.Printf("%+v\n", evt)
	}

}
