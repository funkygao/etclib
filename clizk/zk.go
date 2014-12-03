package clizk

import (
	log "github.com/funkygao/log4go"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

type CliZk struct {
	mu           sync.Mutex
	servers      []string
	timeout      time.Duration
	sessionEvent <-chan zk.Event
	clientRedial chan bool
	client       *zk.Conn
}

func New() *CliZk {
	return &CliZk{}
}

func (this *CliZk) DialTimeout(servers []string, timeout time.Duration) error {
	if this.client != nil {
		log.Warn("zk[%+v] dial while already connected", servers)
		return nil
	}

	client, sessionChan, err := zk.Connect(servers, timeout)
	if err != nil {
		return err
	}

	this.clientRedial = make(chan bool)
	this.client = client
	this.servers = servers
	this.timeout = timeout
	this.sessionEvent = sessionChan

	go this.runWatchdog()

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
	close(this.clientRedial)
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

func (this *CliZk) Children(path string) ([]string, error) {
	keys, _, err := this.client.Children(path)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (this *CliZk) Delete(path string) error {
	return this.client.Delete(path, -1)
}

func (this *CliZk) WatchChildren(path string, ch chan []string) (err error) {
	var (
		children     []string
		watchEvtChan <-chan zk.Event
	)

	children, _, watchEvtChan, err = this.client.ChildrenW(path)
	go func() {
		for {
			select {
			case evt := <-watchEvtChan:
				log.Trace("zk event: {%s %s %s}",
					evt.Path,
					evt.Type.String(), evt.State.String())

				ch <- children

				children, _, watchEvtChan, err = this.client.ChildrenW(path)
			}
		}
	}()
	ch <- children
	return err
}

func (this *CliZk) NodeExistsError(err error) bool {
	return err == zk.ErrNodeExists
}

func (this *CliZk) defaultAcls() []zk.ACL {
	return zk.WorldACL(zk.PermAll)
}

func (this *CliZk) runWatchdog() {
	var (
		evt zk.Event
		err error
	)

	for {
		select {
		case evt = <-this.sessionEvent:
			if evt.Type == zk.EventSession && evt.State == zk.StateExpired {
				this.Close()

				// redial
				err = this.DialTimeout(this.servers, this.timeout)
				if err != nil {
					log.Error("zk[%+v]: %s", this.servers, err)
					return
				} else {
					log.Trace("zk[%+v] redialed ok", this.servers)
				}

				// notify other goroutins
				this.clientRedial <- true
			}
		}
	}

}
