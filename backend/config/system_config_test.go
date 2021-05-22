package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestSystemConfigInterpolation(t *testing.T) {
	configFile := tempFile()
	content := `
[platform]
data_dir: test
nginx_config_dir: %(data_dir)s/dir
`

	err := ioutil.WriteFile(configFile.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	db := tempFile().Name()
	_ = os.Remove(db)
	config, err := NewSystemConfig(configFile.Name())
	assert.Nil(t, err)
	dir, err := config.NginxConfigDir()
	assert.Nil(t, err)
	assert.Equal(t, "test/dir", *dir)
}
