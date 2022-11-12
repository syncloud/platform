package backup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"

	"github.com/stretchr/testify/assert"
)

type DiskUsageStub struct {
	used uint64
}

func (e *DiskUsageStub) Used(_ string) (uint64, error) {
	return e.used, nil
}

type SnapServiceStub struct {
	versionDir string
}

func (s *SnapServiceStub) Stop(_ string) error {
	return nil
}

func (s *SnapServiceStub) Start(_ string) error {
	return nil
}

func (s *SnapServiceStub) RunCmdIfExists(_ model.Snap, cmd string) error {
	if cmd == CreatePreStop {
		backupFile := filepath.Join(s.versionDir, "backup.file")
		if err := os.WriteFile(backupFile, []byte("backup"), 0666); err != nil {
			panic(err)
		}
	}
	return nil
}

type SnapInfoStub struct {
}

func (s *SnapInfoStub) Snap(_ string) (model.Snap, error) {
	return model.Snap{}, nil
}

func TestRemove(t *testing.T) {
	logger := log.Default()

	backupDir := t.TempDir()
	varDir := t.TempDir()
	tmpfn := filepath.Join(backupDir, "tmpfile")
	if err := os.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		panic(err)
	}
	backup := New(backupDir, varDir, cli.New(log.Default()), &DiskUsageStub{100}, &SnapServiceStub{}, &SnapInfoStub{}, logger)
	err := backup.Remove("tmpfile")
	assert.Nil(t, err)
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestBackup(t *testing.T) {
	backupDir := t.TempDir()
	varDir := t.TempDir()
	appDir := filepath.Join(varDir, "test-app")
	_ = os.Mkdir(appDir, 0750)
	versionDir := filepath.Join(appDir, "x1")
	_ = os.Mkdir(versionDir, 0750)
	currentDir := filepath.Join(appDir, "current")
	_ = os.Symlink(versionDir, currentDir)
	commonDir := filepath.Join(appDir, "common")
	_ = os.Mkdir(commonDir, 0750)

	currentFile := filepath.Join(versionDir, "current.file")
	if err := os.WriteFile(currentFile, []byte("current"), 0666); err != nil {
		panic(err)
	}

	commonFile := filepath.Join(commonDir, "common.file")
	if err := os.WriteFile(commonFile, []byte("common"), 0666); err != nil {
		panic(err)
	}

	diskusage := &DiskUsageStub{100}
	logger := log.Default()
	shellExecutor := cli.New(logger)
	snapCli := &SnapServiceStub{versionDir: versionDir}
	snapServer := &SnapInfoStub{}
	backup := New(backupDir+"/non-existent", varDir, shellExecutor, diskusage, snapCli, snapServer, logger)
	backup.Init()
	err := backup.Create("test-app")
	assert.Nil(t, err)
	backups, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(backups), 1)

	err = os.Remove(currentFile)
	assert.Nil(t, err)

	err = os.Remove(commonFile)
	assert.Nil(t, err)

	err = backup.Restore(backups[0].File)
	assert.Nil(t, err)

	currentFileContent, err := os.ReadFile(currentFile)
	assert.Nil(t, err)
	assert.Equal(t, "current", string(currentFileContent))

	backupFileContent, err := os.ReadFile(filepath.Join(versionDir, "backup.file"))
	assert.Nil(t, err)
	assert.Equal(t, "backup", string(backupFileContent))

	commonFileContent, err := os.ReadFile(commonFile)
	assert.Nil(t, err)
	assert.Equal(t, "common", string(commonFileContent))

}
