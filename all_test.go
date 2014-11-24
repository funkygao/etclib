package etclib

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestKeyPath(t *testing.T) {
	project = "dw"
	assert.Equal(t, "/proj/dw/fae/12.12.21.21:1922",
		keyPath("fae", "12.12.21.21:1922"))
	assert.Equal(t, "/proj/dw/node/fae/33.31.12.54:1233",
		keyPath(DIR_NODES, "fae", "33.31.12.54:1233"))
	assert.Equal(t, "/proj/dw/maintain", keyPath("maintain"))
}

func TestNodePath(t *testing.T) {
	project = "dw"
	assert.Equal(t, "/proj/dw/node/fae/12.32.1.5:9001",
		nodePath("fae", "12.32.1.5:9001"))
}

func TestNodeRoot(t *testing.T) {
	project = "dw"
	assert.Equal(t, "/proj/dw/node/actor", nodeRoot(NODE_ACTOR))
}

func BenchmarkKeyPath(b *testing.B) {
	project = "dw"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		keyPath("fae", "12.12.21.21:1922")
	}
}
