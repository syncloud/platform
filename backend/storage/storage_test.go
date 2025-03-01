package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"
)

type StorageExecutorStub struct {
	output string
}

func (e *StorageExecutorStub) CombinedOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

type StorageConfigStub struct {
	diskDir string
}

func (c *StorageConfigStub) DiskLink() string {
	return c.diskDir
}

func TestStorage_ChownRecursive(t *testing.T) {
	storageDir := t.TempDir()
	app1Dir := filepath.Join(storageDir, "app1")
	err := os.MkdirAll(app1Dir, 0777)
	assert.NoError(t, err)
	dataLink := filepath.Join(storageDir, "data")
	output, err := exec.Command("ln", "-s", app1Dir, dataLink).CombinedOutput()
	assert.NoError(t, err, output)

	assert.Nil(t, os.WriteFile(filepath.Join(app1Dir, "1"), []byte(""), 0666))
	assert.Nil(t, os.WriteFile(filepath.Join(app1Dir, "2"), []byte(""), 0666))
	assert.Nil(t, os.WriteFile(filepath.Join(app1Dir, "3"), []byte(""), 0666))

	storage := New(
		&StorageConfigStub{},
		&StorageExecutorStub{},
		log.Default())

	currentUser, err := user.Current()
	assert.Nil(t, err)
	err = storage.ChownRecursive(dataLink, currentUser.Username)

	assert.Nil(t, err)
}

func TestStorage_InitAppStorage(t *testing.T) {
	storageDir := t.TempDir()

	storage := New(
		&StorageConfigStub{diskDir: storageDir},
		&StorageExecutorStub{},
		log.Default())

	path, err := storage.InitAppStorage("app1")

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s/app1", storageDir), path)
}

func TestStorage_InitAppStorageOwner(t *testing.T) {
	storageDir := t.TempDir()

	storage := New(
		&StorageConfigStub{diskDir: storageDir},
		&StorageExecutorStub{},
		log.Default())
	currentUser, err := user.Current()
	assert.Nil(t, err)
	path, err := storage.InitAppStorageOwner("app1", currentUser.Username)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s/app1", storageDir), path)
}
