package clizk

import (
	log "github.com/funkygao/log4go"
	"sort"
)

type Master struct {
	cli *CliZk

	root string // root path in zk
	num  uint16 // how many masters can coexist at the same time
}

func NewMaster(num int, root string) (this *Master) {
	this = &Master{num: uint16(num), root: root}
	return
}

func (this *Master) Join(val string) error {
	return this.cli.Create(this.root, val, FlagServiceSequence)
}

func (this *Master) WatchWho(chMasterVals chan []string) {
	ch := make(chan []string, 10)
	go this.cli.WatchChildren(this.root, ch)

	for {
		select {
		case <-ch:
			children, err := this.cli.Children(this.root)
			if err == nil {
				if len(children) == 0 {
					log.Warn("found no master candidates")
				} else {
					// sort the children
					sort.Sort(sort.StringSlice(children))
					masterVals := make([]string, 0)
					for _, masterNode := range children[:this.num] {
						v, e := this.cli.Get(masterNode)
						if e != nil {
							log.Error("%s", e)
							continue
						}

						masterVals = append(masterVals, v)
					}

					chMasterVals <- masterVals
				}
			} else {
				log.Error("%s", err)
			}
		}
	}

}
