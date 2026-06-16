package auth

import (
	"strings"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
)

func attr(req *ldap.AddRequest, name string) []string {
	for _, a := range req.Attributes {
		if a.Type == name {
			return a.Vals
		}
	}
	return nil
}

func TestUserBuilder_MatchUsersApp(t *testing.T) {
	req := NewUserBuilder().Build("bob", "bob@example.com", 2001, "Password1")
	assert.Equal(t, "cn=bob,ou=users,dc=syncloud,dc=org", req.DN)
	assert.Equal(t, []string{"bob"}, attr(req, "cn"))
	assert.Equal(t, []string{"bob"}, attr(req, "sn"))
	assert.Equal(t, []string{"bob"}, attr(req, "givenName"))
	assert.Equal(t, []string{"bob"}, attr(req, "displayName"))
	assert.Equal(t, []string{"bob"}, attr(req, "uid"))
	assert.Equal(t, []string{"bob@example.com"}, attr(req, "mail"))
	assert.Equal(t, []string{"/home/bob"}, attr(req, "homeDirectory"))
	assert.Equal(t, []string{"/bin/bash"}, attr(req, "loginShell"))
	assert.Equal(t, []string{"2001"}, attr(req, "uidNumber"))
	assert.Equal(t, []string{"2001"}, attr(req, "gidNumber"))
	assert.Contains(t, attr(req, "objectClass"), "posixAccount")
	assert.Contains(t, attr(req, "objectClass"), "inetOrgPerson")
	assert.Contains(t, attr(req, "objectClass"), "person")
	assert.Contains(t, attr(req, "objectClass"), "simpleSecurityObject")
	assert.True(t, strings.HasPrefix(attr(req, "userPassword")[0], "{SSHA}"))
}
