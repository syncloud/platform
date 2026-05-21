//go:build linux

package stability

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	swapFlagPrefer   = 0x8000
	swapFlagPrioMask = 0x7fff
)

func swaponFlags(priority int) int {
	return swapFlagPrefer | (priority & swapFlagPrioMask)
}

func SwaponSyscall(path string, flags int) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall(unix.SYS_SWAPON, uintptr(unsafe.Pointer(p)), uintptr(flags), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func SwapoffSyscall(path string) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall(unix.SYS_SWAPOFF, uintptr(unsafe.Pointer(p)), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func deviceSize(f *os.File, statSize uint64) (uint64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	if fi.Mode()&os.ModeDevice == 0 {
		return statSize, nil
	}
	var n uint64
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), unix.BLKGETSIZE64, uintptr(unsafe.Pointer(&n)))
	if errno != 0 {
		return 0, errno
	}
	return n, nil
}
