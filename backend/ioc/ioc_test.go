package ioc

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"testing"
)

func TestIoC(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	configDb, err := ioutil.TempFile("", "")
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

	backupDir, err := ioutil.TempFile("", "")
	assert.Nil(t, err)

	Init(configDb.Name(), systemConfig.Name(), backupDir.Name(), logger)
}
