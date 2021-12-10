package cert

import (
	"go.uber.org/zap"
	"os/exec"
)

const Subject = "/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud"

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

	output, err := exec.Command("snap",
		"run", "platform.openssl",
		"req",
		"-x509", "-nodes",
		"-newkey", "rsa:2048",
		"-keyout", c.systemConfig.SslKeyFile(),
		"-out", c.systemConfig.SslCertificateFile(),
		"-days", "3650",
		"-subj", Subject).CombinedOutput()
	c.logger.Info("openssl output", zap.String("output", string(output)))
	return err
}
