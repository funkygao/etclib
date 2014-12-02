package etclib

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestFae(t *testing.T) {
	defer Close()
	err := Dial([]string{"127.0.0.1:2181"})
	assert.Equal(t, nil, err)
	BootService("182.11.33.11:9001", SERVICE_FAE)
	endpoints, err := ServiceEndpoints(SERVICE_FAE)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(endpoints))
	t.Logf("%+v", endpoints)

	BootService("182.11.33.12:9001", SERVICE_FAE)
	endpoints, err = ServiceEndpoints(SERVICE_FAE)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(endpoints))
	t.Logf("%+v", endpoints)

	ShutdownService("182.11.33.12:9001", SERVICE_FAE)
	endpoints, err = ServiceEndpoints(SERVICE_FAE)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(endpoints))

	ShutdownService("182.11.33.11:9001", SERVICE_FAE)
	endpoints, err = ServiceEndpoints(SERVICE_FAE)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(endpoints))

}
