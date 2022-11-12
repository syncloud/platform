package config

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemConfigInterpolation(t *testing.T) {
	configFile := tempFile()
	content := `
[platform]
app_dir: test
config_dir: %(app_dir)s/dir
`

	err := os.WriteFile(configFile.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = os.Remove(configFile.Name()) }()
	config := NewSystemConfig(configFile.Name())
	config.Load()
	dir := config.ConfigDir()
	assert.Equal(t, "test/dir", dir)
}
