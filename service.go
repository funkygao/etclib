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

	log.Debug("etcd[%s] connected with %+v", project, servers)
}

func setNodeStatus(nodeType, nodeAddr, status string) {
	if client == nil {
		panic("Call Init before this")
	}

	key := nodePath(nodeType, nodeAddr)
	log.Debug("etcd[%s] node[%s] -> %s", project, key, status)
	_, err := client.Set(key, status, NODE_PING_INTERVAL)
	if err != nil {
		log.Error("etcd[%s] node[%s]: %s", project, key, err.Error())
	}

	go func(key string) {
		// ttl=30，那么应该还是每30s去renew，但提前4s FIXME
		ticker := time.NewTicker(time.Duration(NODE_PING_INTERVAL-4) * time.Second)
		defer ticker.Stop()

		for _ = range ticker.C {
			log.Debug("etcd[%s] node[%s] heartbeat", project, key)

			_, err := client.Set(key, status, NODE_PING_INTERVAL)
			if err != nil {
				log.Error("etcd[%s] node[%s]: %s", project, key, err.Error())
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
	log.Debug("etcd[%s] node[%s] -> shutdown", project, key)
	_, err := client.Delete(key, false)
	if err != nil {
		log.Error("etcd[%s] node[%s]: %s", project, key, err.Error())
	}
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

// who are under maintenance
func MaintainInfo() ([]string, error) {
	resp, err := client.Get(maintainRoot(), false, true)
	if err != nil {
		return nil, err
	}

	info := make([]string, 0)
	for _, node := range resp.Node.Nodes {
		if node.Value == MAINTAIN_EVT_MAINTAIN {
			info = append(info, node.Key)
		}
	}

	return info, nil
}

func WatchFaeNodes() (ch chan NodeEvent) {
	return watchNodes(NODE_FAE)
}

func WatchActorNodes() (ch chan NodeEvent) {
	return watchNodes(NODE_ACTOR)
}

func WatchMaintain() (ch chan MaintainEvent) {
	ch = make(chan MaintainEvent, 10)
	watchChan := make(chan *etcd.Response)

	go func() {
		// TODO add waitIndex to catchup
		_, err := client.Watch(maintainRoot(), 0, true, watchChan, nil)
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		for evt := range watchChan {
			log.Debug("event[%s] %+v", evt.Action, *evt.Node)

			switch evt.Action {
			case "set", "create":
				ch <- MaintainEvent{Key: nodeName(evt.Node.Key),
					Value: MAINTAIN_EVT_MAINTAIN}

			case "delete":
				ch <- MaintainEvent{Key: nodeName(evt.Node.Key),
					Value: MAINTAIN_EVT_UNMAINTAIN}
			}
		}
	}()

	return
}

func watchNodes(nodeType string) (ch chan NodeEvent) {
	ch = make(chan NodeEvent, 10)
	if _, present := nodes[nodeType]; !present {
		nodes[nodeType] = make(map[string]bool)
	}

	watchChan := make(chan *etcd.Response)
	// http long polling, auto reconnect HTTP
	go func() {
		_, err := client.Watch(nodeRoot(nodeType), 0, true, watchChan, nil)
		if err != nil {
			log.Error("node err: %s", err)
		}
	}()

	go func() {
		for evt := range watchChan {
			switch evt.Action {
			case "set", "create":
				if _, present := nodes[nodeType][evt.Node.Key]; !present {
					nodes[nodeType][evt.Node.Key] = true

					ch <- NodeEvent{Addr: nodeName(evt.Node.Key),
						EventType: NODE_EVT_BOOT, NodeType: nodeType}
				}

			case "delete", "expire":
				delete(nodes[nodeType], evt.Node.Key)
				ch <- NodeEvent{Addr: nodeName(evt.Node.Key),
					EventType: NODE_EVT_SHUTDOWN, NodeType: nodeType}
			}

		}
	}()

	return
}
