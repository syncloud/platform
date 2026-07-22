package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsernameValidator_Valid(t *testing.T) {
	validator := NewUsernameValidator()
	for _, username := range []string{
		"bob",
		"ab",
		"user1",
		"first.last",
		"first-last",
		"first_last",
		"a234567890123456789012345678901x",
	} {
		assert.NoError(t, validator.Validate(username), username)
	}
}

func TestUsernameValidator_Invalid(t *testing.T) {
	validator := NewUsernameValidator()
	for _, username := range []string{
		"",
		"a",
		"Bob",
		"bOb",
		"1bob",
		".bob",
		"bob smith",
		"bòb",
		"bob@example.com",
		"cn=admin,ou=users",
		"a2345678901234567890123456789012x",
	} {
		assert.Error(t, validator.Validate(username), username)
	}
}
