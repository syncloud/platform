package backup

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
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
	fmt.Println("stop")
	return nil
}

func (s *SnapServiceStub) Start(_ string) error {
	fmt.Println("start")
	return nil
}

func (s *SnapServiceStub) RunCmdIfExists(_ model.Snap, cmd string) error {
	fmt.Println("run cmd", cmd)
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

func (s *SnapInfoStub) FindInstalled(_ string) (*model.Snap, error) {
	return &model.Snap{}, nil
}

type UserConfigStub struct {
	auto string
	day  int
	hour int
}

func (u *UserConfigStub) GetBackupAuto() string {
	return u.auto
}

func (u *UserConfigStub) SetBackupAuto(auto string) {
	u.auto = auto
}

func (u *UserConfigStub) GetBackupAutoDay() int {
	return u.day
}

func (u *UserConfigStub) SetBackupAutoDay(day int) {
	u.day = day
}

func (u *UserConfigStub) GetBackupAutoHour() int {
	return u.hour
}

func (u *UserConfigStub) SetBackupAutoHour(hour int) {
	u.hour = hour
}

type ProviderStub struct {
	now time.Time
}

func (p ProviderStub) Now() time.Time {
	p.now = p.now.AddDate(0, 0, 1)
	return p.now
}

func TestBackup_Remove(t *testing.T) {
	backupDir := t.TempDir()
	varDir := t.TempDir()
	tmpFile := filepath.Join(backupDir, "tmpfile")
	err := os.WriteFile(tmpFile, []byte(""), 0666)
	assert.NoError(t, err)
	backup := New(
		backupDir,
		varDir,
		cli.New(log.Default()),
		&DiskUsageStub{100},
		&SnapServiceStub{},
		&SnapInfoStub{},
		&UserConfigStub{},
		&ProviderStub{},
		log.Default())
	err = backup.Remove("tmpfile")
	assert.Nil(t, err)
	list, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 0)
}

func TestBackup_Create(t *testing.T) {
	backupDir := t.TempDir()
	varDir := t.TempDir()
	appDir := filepath.Join(varDir, "test-app")
	_ = os.Mkdir(appDir, 0750)
	version := "x1"
	versionDir := filepath.Join(appDir, version)
	_ = os.Mkdir(versionDir, 0750)
	currentDir := filepath.Join(appDir, "current")
	_ = os.Symlink(version, currentDir)
	commonDir := filepath.Join(appDir, "common")
	_ = os.Mkdir(commonDir, 0750)
	socketCommonFile := filepath.Join(commonDir, "web.socket")
	_, err := net.Listen("unix", socketCommonFile)
	assert.Nil(t, err)

	currentFile := filepath.Join(versionDir, "current.file")
	if err := os.WriteFile(currentFile, []byte("current"), 0666); err != nil {
		panic(err)
	}
	socketVersionFile := filepath.Join(versionDir, "web.socket")
	_, err = net.Listen("unix", socketVersionFile)
	assert.Nil(t, err)

	commonFile := filepath.Join(commonDir, "common.file")
	if err := os.WriteFile(commonFile, []byte("common"), 0666); err != nil {
		panic(err)
	}

	app := "test-app"

	err = linux.CreateUser(app)
	assert.NoError(t, err)

	backup := New(
		backupDir+"/non-existent",
		varDir,
		cli.New(log.Default()),
		&DiskUsageStub{100},
		&SnapServiceStub{versionDir: versionDir},
		&SnapInfoStub{},
		&UserConfigStub{},
		&ProviderStub{},
		log.Default())
	backup.Start()
	err = backup.Create(app)
	assert.Nil(t, err)
	backups, err := backup.List()
	assert.Nil(t, err)
	assert.Equal(t, len(backups), 1)

	toDeleteFile := filepath.Join(currentDir, "file.to.delete")
	err = os.WriteFile(toDeleteFile, []byte("test"), 0666)
	assert.NoError(t, err)

	err = os.Remove(currentFile)
	assert.Nil(t, err)

	err = os.Remove(commonFile)
	assert.Nil(t, err)

	err = backup.Restore(backups[0].File)
	assert.Nil(t, err)

	_, err = os.Stat(toDeleteFile)
	assert.ErrorIs(t, err, os.ErrNotExist, "toDeleteFile should not exist")
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

func TestBackup_Auto(t *testing.T) {
	backupDir := t.TempDir()
	varDir := t.TempDir()
	tmpFile := filepath.Join(backupDir, "tmpfile")
	err := os.WriteFile(tmpFile, []byte(""), 0666)
	assert.NoError(t, err)
	backup := New(
		backupDir,
		varDir,
		cli.New(log.Default()),
		&DiskUsageStub{100},
		&SnapServiceStub{},
		&SnapInfoStub{},
		&UserConfigStub{auto: "no", day: 0, hour: 0},
		&ProviderStub{},
		log.Default())

	auto := backup.Auto()
	assert.Equal(t, "no", auto.Auto)
	assert.Equal(t, 0, auto.Day)
	assert.Equal(t, 0, auto.Hour)

	backup.SetAuto(Auto{Auto: "backup", Day: 1, Hour: 2})

	auto = backup.Auto()
	assert.Equal(t, "backup", auto.Auto)
	assert.Equal(t, 1, auto.Day)
	assert.Equal(t, 2, auto.Hour)
}
