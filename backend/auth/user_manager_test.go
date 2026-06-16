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
		NewPasswordValidator(),
		hasher,
		NewEmailResolver(DomainProviderStub{domain: domain}),
		NewUserBuilder(hasher),
	)
}

func TestAddUser_WeakPasswordRejected(t *testing.T) {
	users := newTestUserManager("example.com")
	err := users.AddUser("bob", "weak", "", false)
	assert.NotNil(t, err)
}

func TestAddUser_EmptyUsernameRejected(t *testing.T) {
	users := newTestUserManager("example.com")
	err := users.AddUser("   ", "password", "", false)
	assert.NotNil(t, err)
}
