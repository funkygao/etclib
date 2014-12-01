package etcdcli

import (
	"github.com/funkygao/go-etcd/etcd"
	"time"
)

type CliEtcd struct {
	client  *etcd.Client
	project string
}

func New() *CliEtcd {
	return &CliEtcd{}
}

func (this *CliEtcd) PrepareDirs(dirs ...string) error {
	for _, dir := range dirs {
		if _, err := this.client.CreateDir(dir, 0); err != nil {
			return err
		}
	}
}

func (this *CliEtcd) Close() {
}

func (this *CliEtcd) DialTimeout(servers []string,
	timeout time.Duration, projectName string) error {
	client := etcd.NewClient(servers)
	client.SetConsistency(etcd.STRONG_CONSISTENCY)
	client.SetDialTimeout(timeout)
	if ok := client.SetCluster(servers); !ok {
		return ErrDialFailure
	}

	this.client = client
	this.project = projectName

	return nil
}

func (this *CliEtcd) KeepValue(key, val string, ttl uint64) error {
	_, err := this.client.Set(key, val, ttl)
	if err != nil {
		return err
	}

	go func(key string) {
		// ttl=30，那么应该还是每30s去renew，但提前4s FIXME
		ticker := time.NewTicker(time.Duration(ttl-4) * time.Second)
		defer ticker.Stop()

		for _ = range ticker.C {
			_, err := this.client.Set(key, val, ttl)
			if err != nil {
				return
			}
		}
	}(key)

	return nil
}

func (this *CliEtcd) Delete(key string) error {
	recursive := false
	if _, err := this.client.Delete(key, recursive); err != nil {
		return err
	}

	return nil
}

func (this *CliEtcd) ChildrenKeys(parentKey string) ([]string, error) {
	sort, recursive := false, true
	resp, err := this.client.Get(parentKey, sort, recursive)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for _, node := range resp.Node.Nodes {
		keys = append(keys, node.Key)
	}

	return keys, nil
}
