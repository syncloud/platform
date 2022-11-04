package backup

import (
	"github.com/syncloud/platform/log"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	backupDir := createTempDir()
	varDir := createTempDir()
	defer os.Remove(backupDir)
	defer os.Remove(varDir)
	tmpfn := filepath.Join(backupDir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		panic(err)
	}
	list, err := New(backupDir, varDir, logger).List()
	assert.Nil(t, err)
	assert.Equal(t, list, []File{File{backupDir, "tmpfile"}})
}

func TestRemove(t *testing.T) {
	logger := log.Default()

	backupDir := createTempDir()
	varDir := createTempDir()
	defer os.Remove(backupDir)
	defer os.Remove(varDir)
	tmpfn := filepath.Join(backupDir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		panic(err)
	}
	backup := New(backupDir, varDir, logger)
	err := backup.Remove("tmpfile")
	assert.Nil(t, err)
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestCreate(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	backupDir := createTempDir()
	varDir := createTempDir()
	defer os.Remove(backupDir)
	defer os.Remove(varDir)
	appDir := filepath.Join(varDir, "test-app")
	os.Mkdir(appDir, 0750)
	currentDir := filepath.Join(appDir, "current")
	os.Mkdir(currentDir, 0750)
	commonDir := filepath.Join(appDir, "common")
	os.Mkdir(commonDir, 0750)
	tmpfn := filepath.Join(currentDir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte("*****************"), 0666); err != nil {
		panic(err)
	}

	backup := New(backupDir+"/new", varDir, logger)
	backup.Create("test-app")
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestStart(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	backupDir := createTempDir()
	varDir := createTempDir()
	defer os.Remove(backupDir)
	defer os.Remove(varDir)

	backup := New(backupDir+"/new", varDir, logger)
	backup.Start()
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func createTempDir() string {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		panic(err)
	}
	return dir
}
