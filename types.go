package etclib

type Node struct {
	Addr string
	Boot bool // True if new node intanced added, else an existing node shutdown
}
