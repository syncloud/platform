package ioc

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestIoC(t *testing.T) {
	configDb, err := ioutil.TempFile("", "")
	_ = os.Remove(configDb.Name())
	assert.Nil(t, err)
	systemConfig, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	content := `
[platform]
app_dir: test
data_dir: test
config_dir: test
`
	err = ioutil.WriteFile(systemConfig.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	backupDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.Remove(backupDir)

	varDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.Remove(varDir)

	Init(configDb.Name(), systemConfig.Name(), backupDir, varDir)
}
