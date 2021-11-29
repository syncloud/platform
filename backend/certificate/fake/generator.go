package fake

import (
	"github.com/syncloud/platform/certificate/config"
	"log"
	"os/exec"
)

type Generator struct {
	systemConfig config.GeneratorSystemConfig
}

func New(systemConfig config.GeneratorSystemConfig) *Generator {
	return &Generator{
		systemConfig: systemConfig,
	}
}

func (c *Generator) Generate() error {
	log.Println("generating self signed certificate")

	output, err := exec.Command("snap",
		"run", "platform.openssl",
		"req",
		"-x509", "-nodes",
		"-newkey", "rsa:2048",
		"-keyout", c.systemConfig.SslKeyFile(),
		"-out", c.systemConfig.SslCertificateFile(),
		"-days", "3650",
		"-subj", "/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud").CombinedOutput()
	log.Printf("openssl output: %s", string(output))
	return err
}
