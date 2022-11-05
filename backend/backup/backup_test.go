package backup

import (
	"fmt"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
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

func (e *ExecutorStub) CommandOutput(name string, args ...string) ([]byte, error) {
	if name == "tar" {
		dir := args[3]
		readDir, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		if len(readDir) < 1 {
			return nil, fmt.Errorf("empty dir")
		}
	}
	return []byte(e.output), nil
}

type DiskUsageStub struct {
	used uint64
}

func (e *DiskUsageStub) Used(_ string) (uint64, error) {
	return e.used, nil
}

type SnapServiceStub struct {
}

func (s *SnapServiceStub) Stop(_ string) error {
	return nil
}

func (s *SnapServiceStub) Start(_ string) error {
	return nil
}

func (s *SnapServiceStub) Run(_ string) error {
	return nil
}

type SnapInfoStub struct {
}

func (s *SnapInfoStub) Snap(_ string) (model.Snap, error) {
	return model.Snap{}, nil
}

func TestList(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	backupDir := createTempDir("backup")
	varDir := createTempDir("var")
	defer os.Remove(backupDir)
	defer os.Remove(varDir)
	tmpfn := filepath.Join(backupDir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		panic(err)
	}
	list, err := New(backupDir, varDir, &ExecutorStub{}, &DiskUsageStub{100}, &SnapServiceStub{}, &SnapInfoStub{}, logger).List()
	assert.Nil(t, err)
	assert.Equal(t, list, []File{{backupDir, "tmpfile"}})
}

func TestRemove(t *testing.T) {
	logger := log.Default()

	backupDir := createTempDir("backup")
	varDir := createTempDir("var")
	defer os.Remove(backupDir)
	defer os.Remove(varDir)
	tmpfn := filepath.Join(backupDir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		panic(err)
	}
	backup := New(backupDir, varDir, &ExecutorStub{}, &DiskUsageStub{100}, &SnapServiceStub{}, &SnapInfoStub{}, logger)
	err := backup.Remove("tmpfile")
	assert.Nil(t, err)
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestCreate(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	backupDir := createTempDir("backup")
	varDir := createTempDir("var")
	defer os.Remove(backupDir)
	defer os.Remove(varDir)
	appDir := filepath.Join(varDir, "test-app")
	_ = os.Mkdir(appDir, 0750)
	versionDir := filepath.Join(appDir, "x1")
	_ = os.Mkdir(versionDir, 0750)
	currentDir := filepath.Join(appDir, "current")
	_ = os.Symlink(versionDir, currentDir)
	commonDir := filepath.Join(appDir, "common")
	_ = os.Mkdir(commonDir, 0750)
	tmpfn := filepath.Join(versionDir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte("*****************"), 0666); err != nil {
		panic(err)
	}

	backup := New(backupDir, varDir, &ExecutorStub{}, &DiskUsageStub{100}, &SnapServiceStub{}, &SnapInfoStub{}, logger)
	err = backup.Create("test-app")
	assert.Nil(t, err)
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestStart(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	backupDir := createTempDir("backup")
	varDir := createTempDir("var")
	defer os.Remove(backupDir)
	defer os.Remove(varDir)

	backup := New(backupDir+"/new", varDir, &ExecutorStub{}, &DiskUsageStub{100}, &SnapServiceStub{}, &SnapInfoStub{}, logger)
	backup.Start()
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func createTempDir(pattern string) string {
	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		panic(err)
	}
	return dir
}
