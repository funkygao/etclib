package etclib

import (
	"github.com/funkygao/etclib/clizk"
	"time"
)

func init() {
	store = clizk.New()
}

func Dial(servers []string) error {
	return DialTimeout(servers, DEFAULT_DIAL_TIMEOUT*time.Second)
}

func DialTimeout(servers []string, timeout time.Duration) error {
	if err := store.DialTimeout(servers, timeout); err != nil {
		store.Close()
		return err
	}

	// always create, even if node already exists
	if err := store.Create("/maintain", "",
		clizk.FlagNormal); !store.NodeExistsError(err) {
		return err
	}
	if err := store.Create("/"+SERVICE_FAE, "",
		clizk.FlagNormal); !store.NodeExistsError(err) {
		return err
	}
	if err := store.Create("/"+SERVICE_ACTOR, "",
		clizk.FlagNormal); !store.NodeExistsError(err) {
		return err
	}

	return nil
}

func IsConnected() bool {
	return store.IsConnected()
}

func Close() {
	store.Close()
}

func BootService(addr string, service string) error {
	if err := checkService(service); err != nil {
		return err
	}

	return store.CreateService("/"+service+"/"+addr, "")
}

func ShutdownService(addr string, service string) error {
	if err := checkService(service); err != nil {
		return err
	}

	return store.Delete("/" + service + "/" + addr)
}

func ServiceEndpoints(service string) ([]string, error) {
	if err := checkService(service); err != nil {
		return nil, err
	}

	return store.Children("/" + service)
}

// will block, caller should put it into goroutine
func WatchChildren(path string, ch chan []string) (err error) {
	return store.WatchChildren(path, ch)
}

func WatchService(service string, ch chan []string) (err error) {
	return store.WatchChildren("/"+service, ch)
}

func Create(path string, value string, flags int32) error {
	return store.Create(path, value, flags)
}

func Delete(path string) error {
	return store.Delete(path)
}

func Children(path string) ([]string, error) {
	return store.Children(path)
}

func checkService(service string) error {
	if service != SERVICE_ACTOR && service != SERVICE_FAE {
		return ErrInvalidService
	}

	return nil
}
