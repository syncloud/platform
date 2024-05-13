package hook

import (
	"errors"
	"fmt"
	cp "github.com/otiai10/copy"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage"
	"go.uber.org/zap"
	"os"
	"path"
)

type Install struct {
	storageChecker storage.Checker
	storageLinker  DisksLinker
	systemConfig   SystemConfig
	certGenerator  CertificateGenerator
	ldap           Ldap
	nginx          Nginx
	web            auth.Web
	logDir         string
	logger         *zap.Logger
}

type SystemConfig interface {
	DiskRoot() string
	InternalDiskDir() string
	DiskLink() string
}

type CertificateGenerator interface {
	Generate() error
}

type DisksLinker interface {
	RelinkDisk(link string, target string) error
}

type Ldap interface {
	Init() error
}

type Nginx interface {
	InitConfig() error
}

const (
	App       = "platform"
	AppDir    = "/snap/platform/current"
	DataDir   = "/var/snap/platform/current"
	CommonDir = "/var/snap/platform/common"
)

func NewInstall(
	storageChecker storage.Checker,
	storageLinker DisksLinker,
	systemConfig SystemConfig,
	certGenerator CertificateGenerator,
	ldap Ldap,
	nginx Nginx,
	web auth.Web,
	logger *zap.Logger,
) *Install {
	return &Install{
		storageChecker: storageChecker,
		storageLinker:  storageLinker,
		systemConfig:   systemConfig,
		certGenerator:  certGenerator,
		ldap:           ldap,
		nginx:          nginx,
		web:            web,
		logDir:         path.Join(CommonDir, "log"),
		logger:         logger,
	}
}

func (i *Install) Install() error {
	err := linux.CreateUser(App)
	if err != nil {
		return err
	}

	err = i.InitConfigs()
	if err != nil {
		return err
	}

	err = i.InitDisk()
	if err != nil {
		return err
	}

	err = i.certGenerator.Generate()
	if err != nil {
		return err
	}
	err = i.ldap.Init()
	if err != nil {
		return err
	}
	err = i.nginx.InitConfig()
	if err != nil {
		return err
	}
	err = i.web.InitConfig()
	if err != nil {
		return err
	}
	return nil
}

func (i *Install) PostRefresh() error {

	err := linux.CreateUser(App)
	if err != nil {
		return err
	}

	err = i.InitConfigs()
	if err != nil {
		return err
	}

	err = i.InitDisk()
	if err != nil {
		return err
	}

	err = i.nginx.InitConfig()
	if err != nil {
		return err
	}

	err = i.web.InitConfig()
	if err != nil {
		return err
	}

	err = cli.Remove(fmt.Sprintf("%s/*.log", i.logDir))
	if err != nil {
		return err
	}

	return nil
}

func (i *Install) InitConfigs() error {
	i.logger.Info("init configs")

	dataDirs := []string{
		i.logDir,
		path.Join(DataDir, "nginx"),
		path.Join(DataDir, "openldap"),
		path.Join(DataDir, "openldap-data"),
	}

	for _, dir := range dataDirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err := cp.Copy(path.Join(AppDir, "certs"), "/usr/share/ca-certificates/mozilla")
	if err != nil {
		return err
	}

	return nil
}

func (i *Install) InitDisk() error {
	i.logger.Info("init disk")
	err := createDir(i.systemConfig.DiskRoot())
	if err != nil {
		return err
	}
	err = createDir(i.systemConfig.InternalDiskDir())
	if err != nil {
		return err
	}

	if !i.storageChecker.ExternalDiskLinkExists() {
		err = i.storageLinker.RelinkDisk(i.systemConfig.DiskLink(), i.systemConfig.InternalDiskDir())
		if err != nil {
			return err
		}
	}
	return nil
}

func createDir(dir string) error {
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(dir, 0755)
	}
	return err
}
