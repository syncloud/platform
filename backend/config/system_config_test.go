package config

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemConfigInterpolation(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	content := `
[platform]
app_dir: test
config_dir: %(app_dir)s/dir
`

	err := os.WriteFile(db, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = os.Remove(db) }()
	config := NewSystemConfig(db)
	config.Load()
	dir := config.ConfigDir()
	assert.Equal(t, "test/dir", dir)
}
