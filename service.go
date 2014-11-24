package etclib

import (
	"github.com/coreos/go-etcd/etcd"
	log "github.com/funkygao/log4go"
	"time"
)

func Init(servers []string, projectName string) {
	project = projectName

	client = etcd.NewClient(servers)
	client.SetConsistency(etcd.STRONG_CONSISTENCY)
	client.SetDialTimeout(time.Second * 4)

	// create all root dirs, ignore prevExist err
	var (
		nodeDirs = []string{
			nodeRoot(NODE_FAE),
			nodeRoot(NODE_ACTOR)}
	)
	for _, dir := range nodeDirs {
		client.CreateDir(dir, 0)
	}

	client.CreateDir(keyPath(DIR_MAINTAIN), 0)
}

func bootNode(nodeType, nodeAddr string) {
	if client == nil {
		panic("Call Init before this")
	}

	ticker := time.NewTicker(time.Duration(NODE_PING_INTERVAL-4) * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		_, err := client.Set(keyPath(nodeType, nodeAddr),
			STATUS_ALIVE, NODE_PING_INTERVAL)
		if err != nil {
			log.Error("%s[%s] config: %s", nodeType, nodeAddr, err.Error())
			return
		}
	}
}

func BootFae(addr string) {
	bootNode(NODE_FAE, addr)
}

func BootActor(addr string) {
	bootNode(NODE_ACTOR, addr)
}

func ClusterNodes(nodeType string) {
	key := nodeRoot(nodeType)
	resp, err := client.Get(key, false, true)
}
