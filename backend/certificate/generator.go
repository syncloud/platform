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
	output, err := exec.Command("snap",
		"platform.openssl", "req", "-x509", "-nodes", "-newkey", "rsa:2048",
		"-keyout", "server.rsa.key",
		"-out", "server.rsa.crt",
		"-days", "3650", "-subj", "\"/C=UK/ST=Syncloud/L=Syncloud/O=Syncloud/CN=syncloud\"").CombinedOutput()
	log.Println(output)
	return err
}
