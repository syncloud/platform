package date

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRealProvider_Now(t *testing.T) {
	provider := New()
	assert.Greater(t, provider.Now().Year(), 2000)
}
