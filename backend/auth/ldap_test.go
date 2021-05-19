package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestGenerateChangePasswordCmd(t *testing.T) {
//	assert.Equal(t, GenerateChangePasswordCmd("123123"), "echo \"root:123123\" | chpasswd")
//	assert.Equal(t, GenerateChangePasswordCmd("123\"123"), "echo \"root:123\\\"123\" | chpasswd")
//	assert.Equal(t, GenerateChangePasswordCmd("123$123"), "echo \"root:123\\$123\" | chpasswd")
//}

func TestToLdapDc(t *testing.T) {
	assert.Equal(t, ToLdapDc("user.syncloud.it"), "dc=user,dc=syncloud,dc=it")
}

//func TestMakeSecret(t *testing.T) {
//	assert.True(t, len(makeSecret("password")) > 1)
//}
