package nginx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicConfig_Substitution(t *testing.T) {
	nginx, _, outputDir := newTestNginx(t, "example.com", nil)
	assert.NoError(t, nginx.InitConfig())
	assertGolden(t, path.Join(outputDir, "nginx.conf"), "nginx.example.com.conf")
}
