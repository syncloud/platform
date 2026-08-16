package nginx

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicConfig_Substitution(t *testing.T) {
	nginx, _, outputDir := newTestNginx(t, "example.com", nil)
	assert.NoError(t, nginx.InitConfig())
	assertGolden(t, path.Join(outputDir, "nginx.conf"), "nginx.example.com.conf")
}

func TestPublicConfig_IndexRevalidates(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Equal(t, 2, strings.Count(config, `add_header Cache-Control "no-cache";`))
}

func TestPublicConfig_AssetsImmutable(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Equal(t, 2, strings.Count(config, `add_header Cache-Control "public, max-age=31536000, immutable";`))
}

func TestPublicConfig_MissingAssetIsNotIndex(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Equal(t, 2, strings.Count(config, "try_files $uri =404;"))
}

func TestPublicConfig_AssetsKeepSecurityHeaders(t *testing.T) {
	config := generatePublicConfig(t)
	for _, block := range assetLocations(config) {
		assert.Contains(t, block, "Strict-Transport-Security")
		assert.Contains(t, block, "Access-Control-Allow-Origin")
	}
}

func generatePublicConfig(t *testing.T) string {
	t.Helper()
	nginx, _, outputDir := newTestNginx(t, "example.com", nil)
	assert.NoError(t, nginx.InitConfig())
	content, err := os.ReadFile(path.Join(outputDir, "nginx.conf"))
	assert.NoError(t, err)
	return string(content)
}

func assetLocations(config string) []string {
	var blocks []string
	for _, part := range strings.Split(config, "location /assets/ {")[1:] {
		blocks = append(blocks, strings.SplitN(part, "}", 2)[0])
	}
	return blocks
}
