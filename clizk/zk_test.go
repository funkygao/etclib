package clizk

import (
	"github.com/funkygao/assert"
	"testing"
	"time"
)

func TestChildren(t *testing.T) {
	cli := New()
	defer cli.Close()
	cli.DialTimeout([]string{"127.0.0.1:2181"}, time.Second)
	keys, err := cli.Children("/")
	assert.Equal(t, nil, err)
	t.Logf("%+v", keys)
}

func TestCRUD(t *testing.T) {
	cli := New()
	defer cli.Close()
	cli.DialTimeout([]string{"127.0.0.1:2181"}, time.Second)
	cli.CreateOrUpdate("/zktest", "foo,bar", 0)
	val, err := cli.Get("/zktest")
	assert.Equal(t, "foo,bar", val)
	assert.Equal(t, nil, err)

	err = cli.Set("/zktest", "newval")
	assert.Equal(t, nil, err)
	val, err = cli.Get("/zktest")
	assert.Equal(t, nil, err)
	assert.Equal(t, "newval", val)

	err = cli.Delete("/zktest")
	assert.Equal(t, nil, err)

	_, err = cli.Get("/zktest") // zk: node does not exist
	assert.NotEqual(t, nil, err)

	cli.Close() // close can be called many times
}

func BenchmarkGet(b *testing.B) {
	cli := New()
	defer cli.Close()
	cli.DialTimeout([]string{"127.0.0.1:2181"}, time.Second)
	cli.CreateOrUpdate("/zktest", "foo,bar", 0)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cli.Get("/zktest")
	}

	cli.Close()
}
