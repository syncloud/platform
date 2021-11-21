package selfsigned

import (
	"log"
	"os/exec"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (c *Generator) Generate() error {
	log.Println("generating self signed certificate")

	output, err := exec.Command("snap",
		"run", "platform.openssl",
		"req",
		"-x509", "-nodes",
		"-newkey", "rsa:2048",
		"-keyout", "server.rsa.key",
		"-out", "server.rsa.crt",
		"-days", "3650",
		"-subj", "/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud").CombinedOutput()
	log.Printf("openssl output: %s", string(output))
	return err
}
