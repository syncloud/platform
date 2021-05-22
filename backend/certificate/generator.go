package certificate

import (
	"log"
	"os/exec"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (c *Generator) GenerateSelfSigned() error {
	log.Println("generating self signed certificate")

	output, err := exec.Command("snap",
		"platform.openssl",
		"req",
		"-x509", "-nodes",
		"-newkey", "rsa:2048",
		"-keyout", "server.rsa.key",
		"-out", "server.rsa.crt",
		"-days", "3650",
		"-subj", "/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud").CombinedOutput()
	log.Printf("openssl output:\n%s", string(output))
	return err
}
