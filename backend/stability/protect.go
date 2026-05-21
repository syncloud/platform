package stability

import "strings"

type Protect struct {
	Comms            []string
	CgroupSubstrings []string
}

func DefaultProtect() Protect {
	return Protect{
		Comms: []string{
			"systemd",
			"sshd",
			"init",
			"snapd",
			"dbus-daemon",
			"login",
		},
		CgroupSubstrings: []string{
			"snap.platform.",
			"snap.users.",
			"init.scope",
			"ssh.service",
			"sshd.service",
		},
	}
}

func (p Protect) IsProtected(v Victim) bool {
	for _, c := range p.Comms {
		if v.Comm == c {
			return true
		}
	}
	for _, s := range p.CgroupSubstrings {
		if strings.Contains(v.Cgroup, s) {
			return true
		}
	}
	return false
}
