package cert

import (
	"log"
	"os/exec"
)

const Subject = "/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud"

type Fake struct {
	systemConfig GeneratorSystemConfig
}

type FakeGenerator interface {
	Generate() error
}

func NewFake(systemConfig GeneratorSystemConfig) *Fake {
	return &Fake{
		systemConfig: systemConfig,
	}
}

func (c *Fake) Generate() error {
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
