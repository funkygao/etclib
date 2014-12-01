package clizk

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type CliZk struct {
	client *zk.Conn
}

func New() *CliZk {
	return &CliZk{}
}

func (this *CliZk) PrepareDirs(dirs ...string) {
	acl := zk.WorldACL(zk.PermAll)
	flags := int32(zk.FlagEphemeral)
	for _, dir := range dirs {
		this.client.Create(dir, []byte(""), flags, acl)
	}
}

func (this *CliZk) DialTimeout(servers []string, timeout time.Duration) error {
	client, _, err := zk.Connect(servers, timeout)
	if err != nil {
		return err
	}

	this.client = client

	return nil
}

func (this *CliZk) Close() {
	this.client.Close()
}

func (this *CliZk) ChildrenKeys(parentKey string) ([]string, error) {
	keys, _, err := this.client.Children(parentKey)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (this *CliZk) Delete(key string) error {
	return this.client.Delete(key, 0)
}
