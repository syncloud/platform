package ioc

import (
	"github.com/syncloud/platform/config"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	configDb, err := os.CreateTemp("", "")
	_ = os.Remove(configDb.Name())
	assert.Nil(t, err)
	systemConfig, err := os.CreateTemp("", "")
	assert.Nil(t, err)
	content := `
[platform]
app_dir: test
data_dir: test
config_dir: test
`
	err = os.WriteFile(systemConfig.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	backupDir := t.TempDir()
	varDir := t.TempDir()

	c, err := Init(configDb.Name(), systemConfig.Name(), backupDir, varDir)
	assert.NoError(t, err)
	var conf *config.SystemConfig
	err = c.Resolve(&conf)
	assert.NoError(t, err)
}
