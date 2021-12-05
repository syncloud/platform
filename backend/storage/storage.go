package storage

import (
	"log"
	"os/exec"
)

type Storage struct {
}

const (
	FormatCmd     = "/snap/platform/current/bin/disk_format.sh"
	BootExtendCmd = "/snap/platform/current/bin/boot_extend.sh"
)

func New() *Storage {
	return &Storage{}
}

func (storage *Storage) Format(device string) {
	log.Println("Running storage format: ", FormatCmd, device)
	out, err := exec.Command(FormatCmd, device).CombinedOutput()
	log.Printf("Storage format output %s", out)
	if err != nil {
		log.Printf("Storage format failed: %v", err)
	} else {
		log.Printf("Storage format completed")
	}
}

func (storage *Storage) BootExtend() {
	log.Println("Running storage boot extend: ", BootExtendCmd)
	out, err := exec.Command(BootExtendCmd).CombinedOutput()
	log.Printf("Storage boot extend output %s", out)
	if err != nil {
		log.Printf("Storage boot extend failed: %v", err)
	} else {
		log.Printf("Storage boot extend completed")
	}
}
