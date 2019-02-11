package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatus(t *testing.T) {
	master := NewMaster()

	assert.Equal(t, master.status, "empty")
}
