package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToLdapDc(t *testing.T) {
	assert.Equal(t, ToLdapDc("user.syncloud.it"), "dc=user,dc=syncloud,dc=it")
}

func TestMakeSecret(t *testing.T) {
	secret := makeSecret("syncloud")
	assert.True(t, len(secret) > 1)
}
