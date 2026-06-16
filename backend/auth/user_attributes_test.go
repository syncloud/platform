package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAttributes_MatchUsersApp(t *testing.T) {
	attrs := NewUserAttributes().Build("bob", "bob@example.com", 2001)
	assert.Equal(t, []string{"bob"}, attrs["cn"])
	assert.Equal(t, []string{"bob"}, attrs["sn"])
	assert.Equal(t, []string{"bob"}, attrs["givenName"])
	assert.Equal(t, []string{"bob"}, attrs["displayName"])
	assert.Equal(t, []string{"bob"}, attrs["uid"])
	assert.Equal(t, []string{"bob@example.com"}, attrs["mail"])
	assert.Equal(t, []string{"/home/bob"}, attrs["homeDirectory"])
	assert.Equal(t, []string{"/bin/bash"}, attrs["loginShell"])
	assert.Equal(t, []string{"2001"}, attrs["uidNumber"])
	assert.Equal(t, []string{"2001"}, attrs["gidNumber"])
	assert.Contains(t, attrs["objectClass"], "posixAccount")
	assert.Contains(t, attrs["objectClass"], "inetOrgPerson")
	assert.Contains(t, attrs["objectClass"], "person")
	assert.Contains(t, attrs["objectClass"], "simpleSecurityObject")
}
