package backup

import (
	"archive/tar"
	"os/user"
	"strconv"
)

type ownerMap struct {
	uids map[string]int
	gids map[string]int
}

func newOwnerMap() *ownerMap {
	return &ownerMap{uids: map[string]int{}, gids: map[string]int{}}
}

func (m *ownerMap) Resolve(header *tar.Header) (int, int) {
	return m.uid(header.Uname, header.Uid), m.gid(header.Gname, header.Gid)
}

func (m *ownerMap) uid(name string, archived int) int {
	if name == "" {
		return archived
	}
	if id, ok := m.uids[name]; ok {
		return id
	}
	id := archived
	if u, err := user.Lookup(name); err == nil {
		if parsed, err := strconv.Atoi(u.Uid); err == nil {
			id = parsed
		}
	}
	m.uids[name] = id
	return id
}

func (m *ownerMap) gid(name string, archived int) int {
	if name == "" {
		return archived
	}
	if id, ok := m.gids[name]; ok {
		return id
	}
	id := archived
	if g, err := user.LookupGroup(name); err == nil {
		if parsed, err := strconv.Atoi(g.Gid); err == nil {
			id = parsed
		}
	}
	m.gids[name] = id
	return id
}
