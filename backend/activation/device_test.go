package activation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmail(t *testing.T) {
	username, email := ParseUsername("test@example.com", "domain")
	assert.Equal(t, "test", username)
	assert.Equal(t, "test@example.com", email)
}

func TestNonEmail(t *testing.T) {
	username, email := ParseUsername("test", "domain")
	assert.Equal(t, "test", username)
	assert.Equal(t, "test@domain", email)
}
