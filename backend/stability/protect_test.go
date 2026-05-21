package stability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProtectMatchesByComm(t *testing.T) {
	p := DefaultProtect()
	assert.True(t, p.IsProtected(Victim{Comm: "sshd"}))
	assert.True(t, p.IsProtected(Victim{Comm: "systemd"}))
	assert.False(t, p.IsProtected(Victim{Comm: "photoprism"}))
}

func TestDefaultProtectMatchesByCgroupSubstring(t *testing.T) {
	p := DefaultProtect()
	assert.True(t, p.IsProtected(Victim{Cgroup: "0::/system.slice/snap.platform.backend.service"}))
	assert.True(t, p.IsProtected(Victim{Cgroup: "1:name=systemd:/system.slice/snap.platform.api.service"}))
	assert.True(t, p.IsProtected(Victim{Cgroup: "0::/system.slice/ssh.service"}))
	assert.False(t, p.IsProtected(Victim{Cgroup: "0::/system.slice/snap.photoprism.web.service"}))
}
