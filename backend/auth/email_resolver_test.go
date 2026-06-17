package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailResolver_EmptyDefaultsToDomain(t *testing.T) {
	resolver := NewEmailResolver(DomainProviderStub{domain: "example.com"})
	email, err := resolver.Resolve("bob", "")
	assert.NoError(t, err)
	assert.Equal(t, "bob@example.com", email)
}

func TestEmailResolver_BlankDefaultsToDomain(t *testing.T) {
	resolver := NewEmailResolver(DomainProviderStub{domain: "example.com"})
	email, err := resolver.Resolve("bob", "   ")
	assert.NoError(t, err)
	assert.Equal(t, "bob@example.com", email)
}

func TestEmailResolver_ValidKept(t *testing.T) {
	resolver := NewEmailResolver(DomainProviderStub{domain: "example.com"})
	email, err := resolver.Resolve("bob", "bob@other.org")
	assert.NoError(t, err)
	assert.Equal(t, "bob@other.org", email)
}

func TestEmailResolver_InvalidRejected(t *testing.T) {
	resolver := NewEmailResolver(DomainProviderStub{domain: "example.com"})
	_, err := resolver.Resolve("bob", "not-an-email")
	assert.Error(t, err)
}
