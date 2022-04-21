
package etcdv3

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	defaultConfig := DefaultConfig()
	assert.Equal(t, time.Second*5, defaultConfig.ConnectTimeout)
	assert.Equal(t, false, defaultConfig.BasicAuth)
	assert.Equal(t, []string(nil), defaultConfig.Endpoints)
	assert.Equal(t, false, defaultConfig.Secure)
}

func TestConfigSet(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"localhost"}
	assert.Equal(t, []string{"localhost"}, config.Endpoints)
}
