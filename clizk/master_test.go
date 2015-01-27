package clizk

import (
	"github.com/funkygao/assert"
	"sort"
	"testing"
)

func TestSorting(t *testing.T) {
	s := []string{"abc", "xmy", "def"}
	sort.Sort(sort.StringSlice(s))
	assert.Equal(t, []string{"abc", "def", "xmy"}, s)
}
