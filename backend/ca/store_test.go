package ca

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

func write(t *testing.T, file, content string) {
	assert.NoError(t, os.MkdirAll(filepath.Dir(file), 0755))
	assert.NoError(t, os.WriteFile(file, []byte(content), 0644))
}

func TestStore_Update(t *testing.T) {
	dir := t.TempDir()
	shipped := filepath.Join(dir, "snap", "certs")
	system := filepath.Join(dir, "usr", "mozilla")
	conf := filepath.Join(dir, "ca-certificates.conf")

	write(t, filepath.Join(shipped, "Kept.crt"), "kept new")
	write(t, filepath.Join(shipped, "Added.crt"), "added")
	write(t, filepath.Join(shipped, "Deselected.crt"), "deselected")
	write(t, filepath.Join(system, "Kept.crt"), "kept old")
	write(t, filepath.Join(system, "Removed.crt"), "removed")
	write(t, filepath.Join(system, "Deselected.crt"), "deselected")
	write(t, conf, "# comment\nmozilla/Kept.crt\nmozilla/Removed.crt\n!mozilla/Deselected.crt\nlocal/Custom.crt\n")

	executor := &ExecutorStub{}
	err := NewStore(shipped, system, conf, executor, log.Default()).Update()
	assert.NoError(t, err)

	content, err := os.ReadFile(conf)
	assert.NoError(t, err)
	assert.Equal(t, "# comment\nmozilla/Kept.crt\n!mozilla/Deselected.crt\nlocal/Custom.crt\nmozilla/Added.crt\n", string(content))

	kept, err := os.ReadFile(filepath.Join(system, "Kept.crt"))
	assert.NoError(t, err)
	assert.Equal(t, "kept new", string(kept))
	assert.FileExists(t, filepath.Join(system, "Added.crt"))
	assert.NoFileExists(t, filepath.Join(system, "Removed.crt"))

	assert.Equal(t, []string{"update-ca-certificates --fresh"}, executor.commands)
}

func TestStore_Update_MissingConf(t *testing.T) {
	dir := t.TempDir()
	shipped := filepath.Join(dir, "certs")
	system := filepath.Join(dir, "mozilla")
	conf := filepath.Join(dir, "ca-certificates.conf")
	write(t, filepath.Join(shipped, "A.crt"), "a")

	err := NewStore(shipped, system, conf, &ExecutorStub{}, log.Default()).Update()
	assert.NoError(t, err)

	content, err := os.ReadFile(conf)
	assert.NoError(t, err)
	assert.Equal(t, "mozilla/A.crt\n", string(content))
}

func TestStore_Update_NothingShipped(t *testing.T) {
	dir := t.TempDir()
	shipped := filepath.Join(dir, "certs")
	assert.NoError(t, os.MkdirAll(shipped, 0755))

	executor := &ExecutorStub{}
	err := NewStore(shipped, filepath.Join(dir, "mozilla"), filepath.Join(dir, "conf"), executor, log.Default()).Update()
	assert.Error(t, err)
	assert.Empty(t, executor.commands)
}
