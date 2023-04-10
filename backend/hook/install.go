package hook

import (
	"errors"
	cp "github.com/otiai10/copy"
	"github.com/syncloud/platform/storage"
	"go.uber.org/zap"
	"os"
	"path"
)

type Install struct {
	storageChecker storage.Checker
	storageLinker  DisksLinker
	config         Config
	certGenerator  CertificateGenerator
	ldap           Ldap
	nginx          Nginx
	logger         *zap.Logger
}

type Config interface {
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

func NewInstall(
	storageChecker storage.Checker,
	storageLinker DisksLinker,
	config Config,
	certGenerator CertificateGenerator,
	ldap Ldap,
	nginx Nginx,
	logger *zap.Logger,
) *Install {
	return &Install{
		storageChecker: storageChecker,
		storageLinker:  storageLinker,
		config:         config,
		certGenerator:  certGenerator,
		ldap:           ldap,
		nginx:          nginx,
		logger:         logger,
	}
}

func (i *Install) Run() error {
	err := i.InitConfigs()
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
	return nil
}

func (i *Install) InitConfigs() error {
	//snapDir := "/snap/platform/current"
	dataDir := "/var/snap/platform/current"
	commonDir := "/var/snap/platform/common"
	//slapdConfigDir := path.Join(dataDir, "slapd.d")

	dataDirs := []string{
		path.Join(commonDir, "log"),
		path.Join(dataDir, "nginx"),
		path.Join(dataDir, "openldap"),
		path.Join(dataDir, "openldap-data"),
	}

	for _, dir := range dataDirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err := cp.Copy("/snap/platform/current/certs", "/usr/share/ca-certificates/mozilla")
	if err != nil {
		return err
	}
	return nil
}

func (i *Install) InitDisk() error {
	err := createDir(i.config.DiskRoot())
	if err != nil {
		return err
	}
	err = createDir(i.config.InternalDiskDir())
	if err != nil {
		return err
	}

	if !i.storageChecker.ExternalDiskLinkExists() {
		err = i.storageLinker.RelinkDisk(i.config.DiskLink(), i.config.InternalDiskDir())
		if err != nil {
			return err
		}
	}
	return nil
}

func createDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
