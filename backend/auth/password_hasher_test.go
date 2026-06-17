package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHasher_Hash(t *testing.T) {
	hash := NewPasswordHasher().Hash("syncloud")
	assert.True(t, strings.HasPrefix(hash, "{SSHA}"))
	assert.Greater(t, len(hash), len("{SSHA}"))
}
