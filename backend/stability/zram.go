package stability

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

const (
	zramDevice          = "/dev/zram0"
	zramModuleName      = "zram"
	zramSysBlockDefault = "/sys/block/zram0"
	zramHotAddDefault   = "/sys/class/zram-control/hot_add"
	procSwapsDefault    = "/proc/swaps"
	zramMaxSizeBytes    = uint64(2 * 1024 * 1024 * 1024)
	zramPriority        = 10
	swapMagicV1         = "SWAPSPACE2"
)

type Zram struct {
	sysBlock  string
	hotAdd    string
	procSwaps string
	devPath   string
	mem       *MemInfo
	events    *EventLog
	log       *zap.Logger
}

func NewZram(mem *MemInfo, events *EventLog, log *zap.Logger) *Zram {
	return &Zram{
		sysBlock:  zramSysBlockDefault,
		hotAdd:    zramHotAddDefault,
		procSwaps: procSwapsDefault,
		devPath:   zramDevice,
		mem:       mem,
		events:    events,
		log:       log,
	}
}

func (z *Zram) EnsureConfigured() error {
	snap, err := z.mem.Snapshot()
	if err != nil {
		return fmt.Errorf("zram: meminfo: %w", err)
	}
	on, err := z.alreadyOn()
	if err != nil {
		return fmt.Errorf("zram: check swaps: %w", err)
	}
	if on {
		z.log.Info("zram: already swap-on", zap.String("dev", z.devPath))
		if err := z.disableFileSwaps(); err != nil {
			z.log.Warn("zram: file-swap disable failed", zap.Error(err))
		}
		return nil
	}
	if err := z.ensureDevice(); err != nil {
		return fmt.Errorf("zram: device: %w", err)
	}
	size := z.sizeBytes(snap.TotalKB * 1024)
	if err := z.configureSysfs(size); err != nil {
		return fmt.Errorf("zram: configure: %w", err)
	}
	if err := mkswapInPlace(z.devPath); err != nil {
		return fmt.Errorf("zram: mkswap: %w", err)
	}
	if err := z.swapon(z.devPath, swaponFlags(zramPriority)); err != nil {
		return fmt.Errorf("zram: swapon: %w", err)
	}
	z.log.Info("zram: enabled", zap.Uint64("size_bytes", size), zap.Int("priority", zramPriority))
	if z.events != nil {
		_ = z.events.Append(Event{Kind: EventKindZramEnabled, SizeBytes: size})
	}
	if err := z.disableFileSwaps(); err != nil {
		z.log.Warn("zram: file-swap disable failed", zap.Error(err))
	}
	return nil
}

func fileSwaps(content string) []string {
	var paths []string
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] == "Filename" {
			continue
		}
		if fields[1] != "file" {
			continue
		}
		paths = append(paths, fields[0])
	}
	return paths
}

func (z *Zram) disableFileSwaps() error {
	b, err := os.ReadFile(z.procSwaps)
	if err != nil {
		return err
	}
	for _, path := range fileSwaps(string(b)) {
		if err := z.swapoff(path); err != nil {
			z.log.Warn("zram: swapoff failed", zap.String("path", path), zap.Error(err))
			continue
		}
		z.log.Info("zram: swapoff file swap", zap.String("path", path))
		if z.events != nil {
			_ = z.events.Append(Event{Kind: EventKindSwapoffFile, Path: path})
		}
	}
	return nil
}

func (z *Zram) alreadyOn() (bool, error) {
	b, err := os.ReadFile(z.procSwaps)
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == z.devPath {
			return true, nil
		}
	}
	return false, nil
}

func (z *Zram) ensureDevice() error {
	if _, err := os.Stat(z.sysBlock); err == nil {
		return nil
	}
	if _, err := os.Stat(z.hotAdd); err != nil {
		if err := z.loadModule(zramModuleName); err != nil {
			z.log.Warn("zram: module load failed", zap.Error(err))
		} else {
			z.log.Info("zram: module loaded", zap.String("name", zramModuleName))
			if z.events != nil {
				_ = z.events.Append(Event{Kind: EventKindZramModuleLoad})
			}
		}
	}
	if _, err := os.Stat(z.sysBlock); err == nil {
		return nil
	}
	if _, err := os.Stat(z.hotAdd); err != nil {
		return fmt.Errorf("no zram and no hot_add at %s", z.hotAdd)
	}
	if _, err := os.ReadFile(z.hotAdd); err != nil {
		return err
	}
	return nil
}

func (z *Zram) sizeBytes(totalBytes uint64) uint64 {
	half := totalBytes / 2
	if half > zramMaxSizeBytes {
		return zramMaxSizeBytes
	}
	return half
}

func (z *Zram) configureSysfs(sizeBytes uint64) error {
	if err := os.WriteFile(filepath.Join(z.sysBlock, "comp_algorithm"), []byte("zstd"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(z.sysBlock, "disksize"), []byte(fmt.Sprintf("%d", sizeBytes)), 0644)
}

func mkswapInPlace(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	pageSize := uint64(os.Getpagesize())
	size, err := deviceSize(f, uint64(st.Size()))
	if err != nil {
		return err
	}
	lastPage := uint32(size/pageSize - 1)
	if _, err := f.WriteAt(make([]byte, 1024), 0); err != nil {
		return err
	}
	header := make([]byte, 1024+128)
	binary.LittleEndian.PutUint32(header[1024:], 1)
	binary.LittleEndian.PutUint32(header[1028:], lastPage)
	binary.LittleEndian.PutUint32(header[1032:], 0)
	if _, err := rand.Read(header[1036 : 1036+16]); err != nil {
		return err
	}
	copy(header[1052:], []byte("zram0"))
	if _, err := f.WriteAt(header[1024:], 1024); err != nil {
		return err
	}
	if _, err := f.WriteAt([]byte(swapMagicV1), int64(pageSize-10)); err != nil {
		return err
	}
	return f.Sync()
}

const (
	swapFlagPrefer   = 0x8000
	swapFlagPrioMask = 0x7fff
)

func swaponFlags(priority int) int {
	return swapFlagPrefer | (priority & swapFlagPrioMask)
}

func (z *Zram) loadModule(name string) error {
	out, err := exec.Command("modprobe", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("modprobe %s: %w: %s", name, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (z *Zram) swapon(path string, flags int) error {
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

func (z *Zram) swapoff(path string) error {
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
