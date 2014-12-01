package etclib

import (
	"github.com/funkygao/etclib/clizk"
	//log "github.com/funkygao/log4go"
	"time"
)

func init() {
	client = clizk.New()
}

func Dial(servers []string) error {
	return client.DialTimeout(servers, DIAL_TIMEOUT*time.Second)
}

func Close() {
	client.Close()
}
