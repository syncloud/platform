package auth

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

type UserManager struct {
	ldapClient        *LdapClient
	groups            *GroupManager
	usernameValidator *UsernameValidator
	passwordValidator *PasswordValidator
	passwordHasher    *PasswordHasher
	emailResolver     *EmailResolver
	userBuilder       *UserBuilder
}

func NewUserManager(ldapClient *LdapClient, groups *GroupManager, usernameValidator *UsernameValidator, passwordValidator *PasswordValidator, passwordHasher *PasswordHasher, emailResolver *EmailResolver, userBuilder *UserBuilder) *UserManager {
	return &UserManager{
		ldapClient:        ldapClient,
		groups:            groups,
		usernameValidator: usernameValidator,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
		emailResolver:     emailResolver,
		userBuilder:       userBuilder,
	}
}

func (m *UserManager) AddUser(username string, password string, email string, admin bool) error {
	if err := m.usernameValidator.Validate(username); err != nil {
		return err
	}
	if err := m.passwordValidator.Validate(password); err != nil {
		return err
	}
	resolvedEmail, err := m.emailResolver.Resolve(username, email)
	if err != nil {
		return err
	}
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)

	id, err := m.nextUid(conn)
	if err != nil {
		return err
	}

	if err := conn.Add(m.userBuilder.Build(username, resolvedEmail, id, password)); err != nil {
		return fmt.Errorf("ldap add user: %w", err)
	}
	if admin {
		return m.groups.AddGroupMember(AdminGroup, username)
	}
	return nil
}

func (m *UserManager) nextUid(conn *ldap.Conn) (int, error) {
	searchRequest := ldap.NewSearchRequest(
		UsersDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixAccount)",
		[]string{"uidNumber"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return 0, fmt.Errorf("ldap uid scan: %w", err)
	}
	next := posixIdStart
	for _, entry := range sr.Entries {
		uid, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if err == nil && uid >= next {
			next = uid + 1
		}
	}
	return next, nil
}

func (m *UserManager) SetUserEmail(username string, email string) error {
	resolvedEmail, err := m.emailResolver.Resolve(username, email)
	if err != nil {
		return err
	}
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	modReq := ldap.NewModifyRequest(userDn, nil)
	modReq.Replace("mail", []string{resolvedEmail})
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap set email: %w", err)
	}
	return nil
}

func (m *UserManager) SetPassword(username string, password string) error {
	if err := m.passwordValidator.Validate(password); err != nil {
		return err
	}
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	modReq := ldap.NewModifyRequest(userDn, nil)
	modReq.Replace("userPassword", []string{m.passwordHasher.Hash(password)})
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap set password: %w", err)
	}
	return nil
}

func (m *UserManager) ListUsers() ([]User, error) {
	groups, err := m.groups.ListGroups()
	if err != nil {
		return nil, err
	}
	membership := map[string][]string{}
	for _, group := range groups {
		for _, member := range group.Members {
			membership[member] = append(membership[member], group.Name)
		}
	}

	conn, err := m.ldapClient.Connect()
	if err != nil {
		return nil, err
	}
	defer m.ldapClient.Disconnect(conn)

	searchRequest := ldap.NewSearchRequest(
		UsersDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=inetOrgPerson)",
		[]string{"cn", "mail"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap list users: %w", err)
	}

	users := make([]User, 0, len(sr.Entries))
	for _, entry := range sr.Entries {
		username := entry.GetAttributeValue("cn")
		userGroups := membership[username]
		admin := false
		other := make([]string, 0, len(userGroups))
		for _, group := range userGroups {
			if group == AdminGroup {
				admin = true
			} else {
				other = append(other, group)
			}
		}
		users = append(users, User{
			Username: username,
			Email:    entry.GetAttributeValue("mail"),
			Admin:    admin,
			Groups:   other,
		})
	}
	return users, nil
}

func (m *UserManager) RemoveUser(username string) error {
	conn, err := m.ldapClient.Connect()
	if err != nil {
		return err
	}
	defer m.ldapClient.Disconnect(conn)

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	delReq := ldap.NewDelRequest(userDn, nil)
	if err := conn.Del(delReq); err != nil {
		return fmt.Errorf("ldap delete user: %w", err)
	}
	return nil
}

func (m *UserManager) SetAdmin(username string, admin bool) error {
	if admin {
		return m.groups.AddGroupMember(AdminGroup, username)
	}
	members, err := m.groups.Members(AdminGroup)
	if err != nil {
		return err
	}
	if len(members) <= 1 && slices.Contains(members, username) {
		return fmt.Errorf("cannot remove the last admin")
	}
	return m.groups.RemoveGroupMember(AdminGroup, username)
}

func (m *UserManager) IsAdmin(username string) (bool, error) {
	members, err := m.groups.Members(AdminGroup)
	if err != nil {
		return false, err
	}
	return slices.Contains(members, username), nil
}
