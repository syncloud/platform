package cert

import (
	"fmt"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

const (
	SubjectCountry      = "UK"
	SubjectProvince     = "Syncloud"
	SubjectLocality     = "Syncloud"
	SubjectOrganization = "Syncloud"
	SubjectCommonName   = "syncloud"
)

type Fake struct {
	systemConfig GeneratorSystemConfig
	logger       *zap.Logger
}

type FakeGenerator interface {
	Generate() error
}

func NewFake(systemConfig GeneratorSystemConfig, logger *zap.Logger) *Fake {
	return &Fake{
		systemConfig: systemConfig,
		logger:       logger,
	}
}

func (c *Fake) Generate() error {
	c.logger.Info("generating self signed certificate")

	subject := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/CN=%s", SubjectCountry, SubjectProvince, SubjectLocality, SubjectOrganization, SubjectCommonName)

	output, err := exec.Command("snap",
		"run", "platform.openssl",
		"req",
		"-x509", "-nodes",
		"-newkey", "rsa:2048",
		"-keyout", c.systemConfig.SslKeyFile(),
		"-out", c.systemConfig.SslCertificateFile(),
		"-days", "3650",
		"-subj", subject).CombinedOutput()

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		c.logger.Info(fmt.Sprintf("openssl output: %s", line))
	}
	return err
}
