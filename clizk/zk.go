package clizk

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type CliZk struct {
	client  *zk.Conn
	project string
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

func (this *CliZk) DialTimeout(servers []string,
	timeout time.Duration, projectName string) error {
	client, _, err := zk.Connect(servers, timeout)
	if err != nil {
		return err
	}

	this.client = client
	this.project = projectName

	return nil
}

func (this *CliZk) Close() {
	this.client.Close()
}
