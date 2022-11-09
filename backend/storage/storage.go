package storage

import (
	"log"
	"os/exec"
)

type DiskStorage interface {
	Format(device string)
	BootExtend()
}

type Storage struct {
}

const (
	FormatCmd     = "/snap/platform/current/bin/disk_format.sh"
	BootExtendCmd = "/snap/platform/current/bin/boot_extend.sh"
)

func New() *Storage {
	return &Storage{}
}

func (storage *Storage) Format(device string) error {
	log.Println("Running storage format: ", FormatCmd, device)
	out, err := exec.Command(FormatCmd, device).CombinedOutput()
	log.Printf("Storage format output %s", out)
	return err
}

func (storage *Storage) BootExtend() error {
	log.Println("Running storage boot extend: ", BootExtendCmd)
	out, err := exec.Command(BootExtendCmd).CombinedOutput()
	log.Printf("Storage boot extend output %s", out)
	if err != nil {
		return err
	}
	return nil
}
