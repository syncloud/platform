package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppName(t *testing.T) {
	for _, tc := range []struct {
		serverName string
		expected   string
	}{
		{"matrix.boris.syncloud.it", "matrix"},
		{"x.matrix.boris.syncloud.it", "matrix"},
		{"Matrix.Boris.Syncloud.It.", "matrix"},
		{"boris.syncloud.it", ""},
		{"matrix.other.syncloud.it", ""},
		{"", ""},
	} {
		assert.Equal(t, tc.expected, AppName(tc.serverName, "boris.syncloud.it"), tc.serverName)
	}
}
