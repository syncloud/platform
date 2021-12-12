package backup

import (
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	dir := createTempDir()
	defer os.RemoveAll(dir)
	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		log.Fatal(err)
	}
	list, err := New(dir, logger).List()
	assert.Nil(t, err)
	assert.Equal(t, list, []File{File{dir, "tmpfile"}})
}

func TestRemove(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	dir := createTempDir()
	defer os.RemoveAll(dir)
	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		log.Fatal(err)
	}
	backup := New(dir, logger)
	err = backup.Remove("tmpfile")
	assert.Nil(t, err)
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestCreateBackupDir(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	dir := createTempDir()
	defer os.RemoveAll(dir)
	backup := New(dir+"/new", logger)
	backup.Start()
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func createTempDir() string {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
