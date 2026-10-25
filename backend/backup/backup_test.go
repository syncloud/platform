package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
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

	err = ChownFile(commonFile, app)
	assert.NoError(t, err)

	err = ChownFile(currentFile, app)
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
	err = backup.Start()
	assert.NoError(t, err)

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
	backupFile := filepath.Join(versionDir, "backup.file")
	backupFileContent, err := os.ReadFile(backupFile)
	assert.Nil(t, err)
	assert.Equal(t, "backup", string(backupFileContent))

	commonFileContent, err := os.ReadFile(commonFile)
	assert.Nil(t, err)
	assert.Equal(t, "common", string(commonFileContent))

	owner, err := getOwner(currentDir)
	assert.NoError(t, err)
	assert.Equal(t, app, owner)

	owner, err = getOwner(currentFile)
	assert.NoError(t, err)
	assert.Equal(t, app, owner)

	owner, err = getOwner(commonFile)
	assert.NoError(t, err)
	assert.Equal(t, app, owner)

}

func ChownFile(file, username string) error {
	command := exec.Command("chown",
		fmt.Sprintf("%s:%s", username, username),
		fmt.Sprintf("%s", file),
	)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}
	return nil
}

func getOwner(file string) (string, error) {
	currentFileInfo, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	stat, ok := currentFileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("not a stat_t")
	}
	u, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		return "", err
	}
	return u.Username, nil
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

func TestTarGz_Symlinks(t *testing.T) {
	sourceDir := t.TempDir()
	snapDir := filepath.Join(sourceDir, "common", "app", "snap", "app")
	assert.NoError(t, os.MkdirAll(filepath.Join(snapDir, "227"), 0755))
	assert.NoError(t, os.Symlink("227", filepath.Join(snapDir, "current")))
	assert.NoError(t, os.Symlink("missing", filepath.Join(sourceDir, "common", "dangling")))
	assert.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "current"), 0755))
	assert.NoError(t, os.WriteFile(filepath.Join(sourceDir, "current", "db.sqlite3"), []byte("vault"), 0644))

	archive := filepath.Join(t.TempDir(), "app.tar.gz")
	assert.NoError(t, createTarGz(archive, sourceDir))

	destDir := t.TempDir()
	assert.NoError(t, extractTarGz(archive, destDir))

	content, err := os.ReadFile(filepath.Join(destDir, "current", "db.sqlite3"))
	assert.NoError(t, err)
	assert.Equal(t, "vault", string(content))

	link, err := os.Readlink(filepath.Join(destDir, "common", "app", "snap", "app", "current"))
	assert.NoError(t, err)
	assert.Equal(t, "227", link)

	link, err = os.Readlink(filepath.Join(destDir, "common", "dangling"))
	assert.NoError(t, err)
	assert.Equal(t, "missing", link)
}

func TestBackup_Restore_RefusesArchiveWithoutCurrent(t *testing.T) {
	backupDir := t.TempDir()
	varDir := t.TempDir()
	appDir := filepath.Join(varDir, "test-app")
	versionDir := filepath.Join(appDir, "x1")
	assert.NoError(t, os.MkdirAll(versionDir, 0750))
	assert.NoError(t, os.Symlink("x1", filepath.Join(appDir, "current")))
	commonDir := filepath.Join(appDir, "common")
	assert.NoError(t, os.Mkdir(commonDir, 0750))
	liveCurrentFile := filepath.Join(versionDir, "config.php")
	assert.NoError(t, os.WriteFile(liveCurrentFile, []byte("live"), 0644))
	liveCommonFile := filepath.Join(commonDir, "common.file")
	assert.NoError(t, os.WriteFile(liveCommonFile, []byte("live"), 0644))

	sourceDir := t.TempDir()
	assert.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "common", "log"), 0755))
	assert.NoError(t, createTarGz(filepath.Join(backupDir, "test-app-2026-1021-203240.tar.gz"), sourceDir))

	backup := New(
		backupDir,
		varDir,
		cli.New(log.Default()),
		&DiskUsageStub{100},
		&SnapServiceStub{versionDir: versionDir},
		&SnapInfoStub{},
		&UserConfigStub{},
		&ProviderStub{},
		log.Default())

	err := backup.Restore("test-app-2026-1021-203240.tar.gz")
	assert.ErrorContains(t, err, "incomplete")

	content, err := os.ReadFile(liveCurrentFile)
	assert.NoError(t, err)
	assert.Equal(t, "live", string(content))
	content, err = os.ReadFile(liveCommonFile)
	assert.NoError(t, err)
	assert.Equal(t, "live", string(content))
}

func TestTarGz_RestoresOwnersByName(t *testing.T) {
	app := "test-app"
	assert.NoError(t, linux.CreateUser(app))
	appUser, err := user.Lookup(app)
	assert.NoError(t, err)

	archive := filepath.Join(t.TempDir(), "app.tar.gz")
	f, err := os.Create(archive)
	assert.NoError(t, err)
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	file := func(name, uname string, uid int) {
		assert.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Typeflag: tar.TypeReg, Mode: 0644, Size: 4, Uid: uid, Gid: uid, Uname: uname, Gname: uname}))
		_, err := tw.Write([]byte("data"))
		assert.NoError(t, err)
	}
	file("./current/db.sqlite3", app, 4242)
	file("./current/box/inbox", "root", 4243)
	file("./current/orphan", "no-such-user-on-this-device", 4244)
	assert.NoError(t, tw.WriteHeader(&tar.Header{Name: "./common/link", Typeflag: tar.TypeSymlink, Linkname: "missing", Uid: 4242, Gid: 4242, Uname: app, Gname: app}))
	assert.NoError(t, tw.Close())
	assert.NoError(t, gw.Close())
	assert.NoError(t, f.Close())

	destDir := t.TempDir()
	assert.NoError(t, extractTarGz(archive, destDir))

	owner, err := getOwner(filepath.Join(destDir, "current", "db.sqlite3"))
	assert.NoError(t, err)
	assert.Equal(t, app, owner)

	owner, err = getOwner(filepath.Join(destDir, "current", "box", "inbox"))
	assert.NoError(t, err)
	assert.Equal(t, "root", owner)

	info, err := os.Stat(filepath.Join(destDir, "current", "orphan"))
	assert.NoError(t, err)
	assert.Equal(t, uint32(4244), info.Sys().(*syscall.Stat_t).Uid)

	info, err = os.Lstat(filepath.Join(destDir, "common", "link"))
	assert.NoError(t, err)
	assert.Equal(t, appUser.Uid, fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid))
}
