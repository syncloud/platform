package fake

import (
	"github.com/syncloud/platform/cert"
	"log"
	"os/exec"
)

const Subject = "/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud"

type Generator struct {
	systemConfig cert.GeneratorSystemConfig
}

func New(systemConfig cert.GeneratorSystemConfig) *Generator {
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
		"-subj", Subject).CombinedOutput()
	log.Printf("openssl output: %s", string(output))
	return err
}
