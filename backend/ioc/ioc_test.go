package ioc

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestIoC(t *testing.T) {
	configDb, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	systemConfig, err := ioutil.TempFile("", "")
	assert.Nil(t, err)

	Init(configDb.Name(), systemConfig.Name())
	//Resolve()
}
