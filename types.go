package etclib

type NodeEvent struct {
	Addr      string // e,g. 12.25.123.55:9002
	NodeType  string
	EventType string
}

type MaintainEvent struct {
	Key   string // global | kingdom_N
	Value string // maintain | unmaintain
}
