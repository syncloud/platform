package storage

import (
	"log"
	"os/exec"
)

type Storage struct {
}

const (
	STORAGE_FORMAT_CMD      = "/snap/platform/current/bin/disk_format.sh"
	STORAGE_BOOT_EXTEND_CMD = "/snap/platform/current/bin/boot_extend.sh"
)

func New() *Storage {
	return &Storage{}
}

func (storage *Storage) Format(device string) {
	log.Println("Running storage format: ", STORAGE_FORMAT_CMD, device)
	out, err := exec.Command(STORAGE_FORMAT_CMD, device).CombinedOutput()
	log.Printf("Storage format output %s", out)
	if err != nil {
		log.Printf("Storage format failed: %v", err)
	} else {
		log.Printf("Storage format completed")
	}
}

func (storage *Storage) BootExtend() {
	log.Println("Running storage boot extend: ", STORAGE_BOOT_EXTEND_CMD)
	out, err := exec.Command(STORAGE_BOOT_EXTEND_CMD).CombinedOutput()
	log.Printf("Storage boot extend output %s", out)
	if err != nil {
		log.Printf("Storage boot extend failed: %v", err)
	} else {
		log.Printf("Storage boot extend completed")
	}
}
