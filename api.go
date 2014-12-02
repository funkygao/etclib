package etclib

import (
	"github.com/funkygao/etclib/clizk"
	"time"
)

func init() {
	store = clizk.New()
}

func Dial(servers []string) error {
	if err := store.DialTimeout(servers, DEFAULT_DIAL_TIMEOUT*time.Second); err != nil {
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

func Set(path, value string) error {
	return store.Set(path, value)
}

func Delete(path string) error {
	return store.Delete(path)
}

func Get(path string) (string, error) {
	return store.Get(path)
}

func checkService(service string) error {
	if service != SERVICE_ACTOR && service != SERVICE_FAE {
		return ErrInvalidService
	}

	return nil
}
