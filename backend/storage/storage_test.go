package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"os"
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

func TestStorage_ChownRecursive_LessThanLimit(t *testing.T) {
	storageDir := t.TempDir()
	assert.Nil(t, os.WriteFile(filepath.Join(storageDir, "1"), []byte(""), 0666))
	storage := New(
		&StorageConfigStub{},
		&StorageExecutorStub{},
		2,
		log.Default())

	currentUser, err := user.Current()
	assert.Nil(t, err)
	changed, err := storage.ChownRecursive(storageDir, currentUser.Username)

	assert.Nil(t, err)
	assert.True(t, changed)
}

func TestStorage_ChownRecursive_MoreThanLimit(t *testing.T) {
	storageDir := t.TempDir()
	assert.Nil(t, os.WriteFile(filepath.Join(storageDir, "1"), []byte(""), 0666))
	assert.Nil(t, os.WriteFile(filepath.Join(storageDir, "2"), []byte(""), 0666))
	assert.Nil(t, os.WriteFile(filepath.Join(storageDir, "3"), []byte(""), 0666))

	storage := New(
		&StorageConfigStub{},
		&StorageExecutorStub{},
		2,
		log.Default())

	currentUser, err := user.Current()
	assert.Nil(t, err)
	changed, err := storage.ChownRecursive(storageDir, currentUser.Username)

	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestStorage_InitAppStorage(t *testing.T) {
	storageDir := t.TempDir()

	storage := New(
		&StorageConfigStub{diskDir: storageDir},
		&StorageExecutorStub{},
		2,
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
		2,
		log.Default())
	currentUser, err := user.Current()
	assert.Nil(t, err)
	path, err := storage.InitAppStorageOwner("app1", currentUser.Username)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s/app1", storageDir), path)
}
