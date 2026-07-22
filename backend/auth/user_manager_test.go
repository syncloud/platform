package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestUserManager(domain string) *UserManager {
	ldapClient := NewLdapClient()
	hasher := NewPasswordHasher()
	return NewUserManager(
		ldapClient,
		NewGroupManager(ldapClient),
		NewUsernameValidator(),
		NewPasswordValidator(),
		hasher,
		NewEmailResolver(DomainProviderStub{domain: domain}),
		NewUserBuilder(hasher),
	)
}

func TestAddUser_WeakPasswordRejected(t *testing.T) {
	users := newTestUserManager("example.com")
	err := users.AddUser("bob", "weak", "", false)
	assert.Error(t, err)
}

func TestAddUser_EmptyUsernameRejected(t *testing.T) {
	users := newTestUserManager("example.com")
	err := users.AddUser("   ", "password1", "", false)
	assert.Error(t, err)
}

func TestAddUser_InvalidUsernameRejected(t *testing.T) {
	users := newTestUserManager("example.com")
	err := users.AddUser("Bob", "password1", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lowercase")
}
