package clizk

import (
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

type CliZk struct {
	mu     sync.Mutex
	client *zk.Conn
}

func New() *CliZk {
	return &CliZk{}
}

func (this *CliZk) PrepareDirs(dirs ...string) {
	flags := int32(zk.FlagEphemeral)
	for _, dir := range dirs {
		this.client.Create(dir, []byte(""), flags, this.defaultAcls())
	}
}

func (this *CliZk) defaultAcls() []zk.ACL {
	return zk.WorldACL(zk.PermAll)
}

// servers item is like ip:port
func (this *CliZk) DialTimeout(servers []string, timeout time.Duration) error {
	client, _, err := zk.Connect(servers, timeout)
	if err != nil {
		return err
	}

	this.client = client

	return nil
}

// reentrant safe
func (this *CliZk) Close() {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.client == nil {
		return
	}

	this.client.Close()
	this.client = nil
}

// not recursive
// flags maybe 0(persistent), zk.FlagEphemeral, FlagSequence
func (this *CliZk) Create(path string, value string, flags int32) error {
	_, err := this.client.Create(path, []byte(value), flags, this.defaultAcls())
	return err
}

func (this *CliZk) CreateOrUpdate(path, value string, flags int32) error {
	err := this.Create(path, value, flags)
	if err == zk.ErrNodeExists {
		err = nil
	}
	return err
}

func (this *CliZk) CreateService(path, value string) error {
	return this.Create(path, value, zk.FlagEphemeral)
}

func (this *CliZk) Exists(path string) (bool, error) {
	r, _, err := this.client.Exists(path)
	return r, err
}

func (this *CliZk) Get(path string) (string, error) {
	val, _, err := this.client.Get(path)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (this *CliZk) Set(path, value string) error {
	_, err := this.client.Set(path, []byte(value), -1)
	return err
}

func (this *CliZk) Children(parentKey string) ([]string, error) {
	keys, _, err := this.client.Children(parentKey)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (this *CliZk) Delete(key string) error {
	return this.client.Delete(key, -1)
}
