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

func setNodeStatus(nodeType, nodeAddr, status string) {
	if client == nil {
		panic("Call Init before this")
	}

	key := nodePath(nodeType, nodeAddr)
	log.Debug("node[%s]: %s", key, status)
	client.Set(key, status, NODE_PING_INTERVAL)

	go func(key string) {
		ticker := time.NewTicker(time.Duration(NODE_PING_INTERVAL-4) * time.Second)
		defer ticker.Stop()

		for _ = range ticker.C {
			_, err := client.Set(key, status, NODE_PING_INTERVAL)
			if err != nil {
				log.Error("node[%s]: %s", key, err.Error())
				return
			}
		}
	}(key)
}

func deleteNode(nodeType, nodeAddr string) {
	if client == nil {
		panic("Call Init before this")
	}

	key := nodePath(nodeType, nodeAddr)
	log.Debug("node[%s]: deleted", key)
	client.Delete(key, false)
}

func BootFae(addr string) {
	setNodeStatus(NODE_FAE, addr, STATUS_ALIVE)
}

func ShutdownFae(addr string) {
	deleteNode(NODE_FAE, addr)
}

func BootActor(addr string) {
	setNodeStatus(NODE_ACTOR, addr, STATUS_ALIVE)
}

func ShutdownActor(addr string) {
	deleteNode(NODE_ACTOR, addr)
}

func ClusterNodes(nodeType string) ([]string, error) {
	key := nodeRoot(nodeType)
	resp, err := client.Get(key, false, true)
	if err != nil {
		return nil, err
	}

	nodes := make([]string, 0)
	for _, node := range resp.Node.Nodes {
		nodes = append(nodes, nodeName(node.Key))
	}

	return nodes, nil
}

func WatchFaeNodes() (ch chan NodeEvent) {
	return watchNodes(NODE_FAE)
}

func WatchActorNodes() (ch chan NodeEvent) {
	return watchNodes(NODE_ACTOR)
}

func watchNodes(nodeType string) (ch chan NodeEvent) {
	ch = make(chan NodeEvent, 10)

	watchChan := make(chan *etcd.Response)
	go client.Watch(nodeRoot(nodeType), 0, true, watchChan, nil)

	go func() {
		var evtType string
		for evt := range watchChan {
			log.Debug("event[%s] %+v", evt.Action, *evt.Node)

			switch evt.Action {
			case "set", "create":
				evtType = NODE_EVT_BOOT

			case "delete":
				evtType = NODE_EVT_SHUTDOWN
			}

			ch <- NodeEvent{Addr: nodeName(evt.Node.Key), EventType: evtType}
		}
	}()

	return
}
