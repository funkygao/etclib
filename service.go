package etclib

import (
	"github.com/coreos/go-etcd/etcd"
	//log "github.com/funkygao/log4go"
	"time"
)

func Init(servers []string, projectName string) {
	project = projectName

	client = etcd.NewClient(servers)
	client.SetConsistency(etcd.STRONG_CONSISTENCY)
	client.SetDialTimeout(time.Second * 4)

	// create all root dirs, ignore prevExist err
	for _, dir := range rootDirs {
		client.CreateDir(dir, 0)
	}

}

func BootFae(addr string) (err error) {
	go func() {
		ticker := time.NewTicker(time.Duration(NODE_PING_INTERVAL))
		defer ticker.Stop()

		for _ = range ticker.C {

			_, err = client.Set(keyPath(DIR_FAE, addr),
				STATUS_ALIVE, NODE_PING_INTERVAL)
			if err != nil {
				return
			}
		}

	}()

	return
}

func BootActor(addr string) (err error) {
	if client == nil {
		panic("")
	}

	return
}

func ClusterServices() {
	client.Get("/", false, true)
}
