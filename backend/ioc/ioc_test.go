package ioc

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIoC(t *testing.T) {
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

	Init(configDb.Name(), systemConfig.Name(), backupDir, varDir)
}
