package etclib

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestKeyPath(t *testing.T) {
	project = "dw"
	assert.Equal(t, "/proj/dw/fae/12.12.21.21:1922", keyPath("fae", "12.12.21.21:1922"))
}

func BenchmarkKeyPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		keyPath("fae", "12.12.21.21:1922")
	}
}
