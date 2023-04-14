package identification

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestMac(t *testing.T) {
	mac, err := GetMac()
	assert.Nil(t, err)
	assert.NotEmpty(t, mac)
}

func TestParser_Id(t *testing.T) {
	tempDir := t.TempDir()
	idFile := path.Join(tempDir, "id")
	err := os.WriteFile(idFile, []byte(`
[id]
name=name
title=title
`), 0644)
	assert.NoError(t, err)
	parser := &Parser{filename: idFile}
	id, err := parser.Id()
	assert.NoError(t, err)
	assert.Equal(t, "name", id.Name)
	assert.Equal(t, "title", id.Title)
}

func TestParser_Id_NoIdSection(t *testing.T) {
	tempDir := t.TempDir()
	idFile := path.Join(tempDir, "id")
	err := os.WriteFile(idFile, []byte(`
name=name
title=title
`), 0644)
	assert.NoError(t, err)
	parser := &Parser{filename: idFile}
	id, err := parser.Id()
	assert.NoError(t, err)
	assert.Equal(t, "unknown", id.Name)
	assert.Equal(t, "Unknown", id.Title)
}
