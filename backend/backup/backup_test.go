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

type ExecutorStub struct {
	output string
}

func (e *ExecutorStub) CommandOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

type DiskUsageStub struct {
	used uint64
}

func (e *DiskUsageStub) Used(_ string) (uint64, error) {
	return e.used, nil
}

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
	list, err := New(backupDir, varDir, &ExecutorStub{}, &DiskUsageStub{100}, logger).List()
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
	backup := New(backupDir, varDir, &ExecutorStub{}, &DiskUsageStub{100}, logger)
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

	backup := New(backupDir+"/new", varDir, &ExecutorStub{}, &DiskUsageStub{100}, logger)
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

	backup := New(backupDir+"/new", varDir, &ExecutorStub{}, &DiskUsageStub{100}, logger)
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
