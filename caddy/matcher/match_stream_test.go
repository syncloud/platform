package matcher

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/mholt/caddy-l4/layer4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func match(t *testing.T, root, serverName string) (bool, string) {
	ctx, cancel := caddy.NewContext(caddy.Context{Context: context.Background()})
	defer cancel()
	matcher := &MatchStream{Domain: "example.com", Root: root}
	require.NoError(t, matcher.Provision(ctx))

	in, out := net.Pipe()
	defer func() {
		_ = out.Close()
		_ = in.Close()
	}()
	go func() {
		_ = tls.Client(in, &tls.Config{ServerName: serverName, InsecureSkipVerify: true}).Handshake()
	}()

	cx := layer4.WrapConnection(out, []byte{}, zap.NewNop())
	matched, err := matcher.Match(cx)
	assert.NoError(t, err)
	app, _ := cx.Context.Value(caddy.ReplacerCtxKey).(*caddy.Replacer).GetString(AppPlaceholder)
	return matched, app
}

func streamSocket(t *testing.T, root, app string) {
	dir := filepath.Join(root, app, "common")
	require.NoError(t, os.MkdirAll(dir, 0755))
	listener, err := net.Listen("unix", filepath.Join(dir, "stream.socket"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })
}

func TestMatchStream_AppWithStreamSocket(t *testing.T) {
	root := t.TempDir()
	streamSocket(t, root, "matrix")
	matched, app := match(t, root, "matrix.example.com")
	assert.True(t, matched)
	assert.Equal(t, "matrix", app)
}

func TestMatchStream_SubdomainOfApp(t *testing.T) {
	root := t.TempDir()
	streamSocket(t, root, "matrix")
	matched, app := match(t, root, "x.matrix.example.com")
	assert.True(t, matched)
	assert.Equal(t, "matrix", app)
}

func TestMatchStream_AppWithoutStreamSocket(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "files", "common"), 0755))
	matched, _ := match(t, root, "files.example.com")
	assert.False(t, matched)
}

func TestMatchStream_RegularFileIsNotASocket(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "files", "common")
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "stream.socket"), nil, 0644))
	matched, _ := match(t, root, "files.example.com")
	assert.False(t, matched)
}

func TestMatchStream_DeviceDomainAndForeignDomain(t *testing.T) {
	root := t.TempDir()
	streamSocket(t, root, "matrix")
	for _, serverName := range []string{"example.com", "matrix.other.com"} {
		matched, _ := match(t, root, serverName)
		assert.False(t, matched, serverName)
	}
}
