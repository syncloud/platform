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
app_dir: test
config_dir: %(app_dir)s/dir
`

	err := ioutil.WriteFile(configFile.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = os.Remove(configFile.Name()) }()
	config, err := NewSystemConfig(configFile.Name())
	assert.Nil(t, err)
	dir := config.ConfigDir()
	assert.Equal(t, "test/dir", dir)
}
