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
	client       *zk.Conn
	quit         chan bool
}

func New() *CliZk {
	return &CliZk{quit: make(chan bool)}
}

func (this *CliZk) DialTimeout(servers []string,
	timeout time.Duration) (err error) {
	if this.client != nil {
		log.Warn("zk[%+v] dial while already connected", servers)
		return nil
	}

	// the client is in StateDisconnected
	this.client, this.sessionEvent, err = zk.Connect(servers, timeout)
	if err != nil {
		return
	}

	this.servers = servers
	this.timeout = timeout

	return
}

func (this *CliZk) IsConnected() bool {
	return this.client != nil
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
				log.Trace("zk event: watch{%s %s %s}",
					evt.Path,
					evt.Type.String(), evt.State.String())

				if evt.Type == zk.EventNodeChildrenChanged {
					ch <- children
					children, _, watchEvtChan, err = this.client.ChildrenW(path)
					log.Trace("zk[%+v] renew children watch", this.servers)
				}

			case evt := <-this.sessionEvent:
				// StateDisconnected vs StateExpired
				// StateDisconnected happens when client can’t connect to the server
				// server is down, then come back, client get StateSyncConnected and
				// everything is back to normal automatically
				// StateExpired happens when client fails to ping server within session
				// timeout period. Server thinks the client is dead and has deleted
				// all the ephemeral nodes and watchers created by the client. It's
				// client's job to reconnect and recreate ephemeral and watchers
				if evt.Type == zk.EventSession && evt.State == zk.StateExpired {
					//
					log.Trace("zk event: session{%s %s}",
						evt.Type.String(), evt.State.String())

					this.Close()

					// redial
					for {
						err = this.DialTimeout(this.servers, this.timeout)
						if err != nil {
							log.Error("zk[%+v]: %s", this.servers, err)
							this.Close()
						} else {
							log.Trace("zk[%+v] redialed ok", this.servers)
							children, _, watchEvtChan, err = this.client.ChildrenW(path)
							break
						}

						time.Sleep(time.Second)
					}
				}

			case <-this.quit: // TODO not used yet
				return
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
