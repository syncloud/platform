package auth

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

var groupNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type GroupManager struct {
	ldapClient *LdapClient
}

func NewGroupManager(ldapClient *LdapClient) *GroupManager {
	return &GroupManager{ldapClient: ldapClient}
}

func (m *GroupManager) ListGroups() ([]Group, error) {
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return nil, err
	}
	defer m.ldapClient.Disconnect(conn)
	return m.listGroups(conn)
}

func (m *GroupManager) listGroups(conn *ldap.Conn) ([]Group, error) {
	searchRequest := ldap.NewSearchRequest(
		GroupsDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"cn", "memberUid", "gidNumber"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap list groups: %w", err)
	}
	groups := make([]Group, 0, len(sr.Entries))
	for _, entry := range sr.Entries {
		groups = append(groups, Group{
			Name:    entry.GetAttributeValue("cn"),
			Members: entry.GetAttributeValues("memberUid"),
		})
	}
	return groups, nil
}

func (m *GroupManager) AddGroup(name string) error {
	if !groupNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid group name: %s", name)
	}
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)

	gid, err := m.nextGid(conn)
	if err != nil {
		return err
	}

	groupDn := fmt.Sprintf("cn=%s,%s", name, GroupsDn)
	addReq := ldap.NewAddRequest(groupDn, nil)
	addReq.Attribute("objectClass", []string{"posixGroup", "top"})
	addReq.Attribute("cn", []string{name})
	addReq.Attribute("gidNumber", []string{strconv.Itoa(gid)})
	if err := conn.Add(addReq); err != nil {
		return fmt.Errorf("ldap add group: %w", err)
	}
	return nil
}

func (m *GroupManager) nextGid(conn *ldap.Conn) (int, error) {
	searchRequest := ldap.NewSearchRequest(
		GroupsDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"gidNumber"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return 0, fmt.Errorf("ldap gid scan: %w", err)
	}
	next := posixIdStart
	for _, entry := range sr.Entries {
		gid, err := strconv.Atoi(entry.GetAttributeValue("gidNumber"))
		if err == nil && gid >= next {
			next = gid + 1
		}
	}
	return next, nil
}

func (m *GroupManager) RemoveGroup(name string) error {
	if name == AdminGroup {
		return fmt.Errorf("cannot remove admin group")
	}
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)

	groupDn := fmt.Sprintf("cn=%s,%s", name, GroupsDn)
	if err := conn.Del(ldap.NewDelRequest(groupDn, nil)); err != nil {
		return fmt.Errorf("ldap remove group: %w", err)
	}
	return nil
}

func (m *GroupManager) AddGroupMember(group string, username string) error {
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)
	return m.addMember(conn, group, username)
}

func (m *GroupManager) RemoveGroupMember(group string, username string) error {
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)
	return m.removeMember(conn, group, username)
}

func (m *GroupManager) Members(group string) ([]string, error) {
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return nil, err
	}
	defer m.ldapClient.Disconnect(conn)
	return m.groupMembers(conn, group)
}

func (m *GroupManager) groupMembers(conn *ldap.Conn, group string) ([]string, error) {
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("cn=%s,%s", group, GroupsDn),
		ldap.ScopeBaseObject, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"memberUid"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap group members: %w", err)
	}
	if len(sr.Entries) < 1 {
		return nil, fmt.Errorf("group not found: %s", group)
	}
	return sr.Entries[0].GetAttributeValues("memberUid"), nil
}

func (m *GroupManager) addMember(conn *ldap.Conn, group string, username string) error {
	members, err := m.groupMembers(conn, group)
	if err != nil {
		return err
	}
	if slices.Contains(members, username) {
		return nil
	}
	modReq := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,%s", group, GroupsDn), nil)
	modReq.Add("memberUid", []string{username})
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap add group member: %w", err)
	}
	return nil
}

func (m *GroupManager) removeMember(conn *ldap.Conn, group string, username string) error {
	members, err := m.groupMembers(conn, group)
	if err != nil {
		return err
	}
	if !slices.Contains(members, username) {
		return nil
	}
	modReq := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,%s", group, GroupsDn), nil)
	modReq.Delete("memberUid", []string{username})
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap remove group member: %w", err)
	}
	return nil
}
