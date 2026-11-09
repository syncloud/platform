package kernel

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestPortResolver_Port(t *testing.T) {
	root := t.TempDir()
	devices := map[string]string{
		"sda":   "devices/platform/soc/fc000000.pcie/pci0000:00/0000:00:00.0/0000:01:00.0/ata1/host0/target0:0:0/0:0:0:0/block/sda",
		"sdb":   "devices/platform/soc/ff500000.usb/xhci-hcd.0.auto/usb2/2-1/2-1:1.0/host2/target2:0:0/2:0:0:0/block/sdb",
		"loop0": "devices/virtual/block/loop0",
	}
	assert.Nil(t, os.MkdirAll(filepath.Join(root, "block"), 0755))
	for name, target := range devices {
		assert.Nil(t, os.MkdirAll(filepath.Join(root, target), 0755))
		assert.Nil(t, os.Symlink(filepath.Join("..", target), filepath.Join(root, "block", name)))
	}
	resolver := NewPortResolver(root)
	assert.Equal(t, "ata1", resolver.Port("/dev/sda"))
	assert.Equal(t, "usb2", resolver.Port("/dev/sdb"))
	assert.Equal(t, "", resolver.Port("/dev/loop0"))
	assert.Equal(t, "", resolver.Port("/dev/sdz"))
}
