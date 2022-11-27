package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse_Simple(t *testing.T) {
	file, err := Parse("/data", "app-2001-020304-050607.tar.gz")
	assert.Nil(t, err)
	assert.Equal(t, file.Path, "/data")
	assert.Equal(t, file.File, "app-2001-020304-050607.tar.gz")
	assert.Equal(t, file.App, "app")
	assert.Equal(t, file.FullName, "/data/app-2001-020304-050607.tar.gz")
}

func TestParse_AppWithDash(t *testing.T) {
	file, err := Parse("/data", "app-name-2001-020304-050607.tar.gz")
	assert.Nil(t, err)
	assert.Equal(t, file.Path, "/data")
	assert.Equal(t, file.File, "app-name-2001-020304-050607.tar.gz")
	assert.Equal(t, file.App, "app-name")
	assert.Equal(t, file.FullName, "/data/app-name-2001-020304-050607.tar.gz")
}

func TestParse_Wrong(t *testing.T) {
	_, err := Parse("/data", "app-name-2001_020304-050607.tar.gz")
	assert.NotNil(t, err)
}
