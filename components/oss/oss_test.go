package oss

import (
	"github.com/stretchr/testify/assert"
	"langgo/core"
	"testing"
)

func TestOss(t *testing.T) {
	core.EnvName = core.Development
	core.LoadConfigurationFile("xxx.yml")
	i := Instance{}
	i.Load()
	object, err := GetObject("xxx")
	assert.NoError(t, err)
	assert.NotNil(t, object)
}
