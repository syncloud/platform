package matcher

import (
	"os"
	"path/filepath"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/mholt/caddy-l4/layer4"
	"github.com/mholt/caddy-l4/modules/l4tls"
)

const AppPlaceholder = "l4.syncloud.app"

func init() {
	caddy.RegisterModule(&MatchStream{})
}

type MatchStream struct {
	Domain string `json:"domain,omitempty"`
	Root   string `json:"root,omitempty"`

	tls *l4tls.MatchTLS
}

func (*MatchStream) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.matchers.syncloud_stream",
		New: func() caddy.Module { return new(MatchStream) },
	}
}

func (m *MatchStream) Provision(ctx caddy.Context) error {
	if m.Root == "" {
		m.Root = "/var/snap"
	}
	m.tls = &l4tls.MatchTLS{}
	return m.tls.Provision(ctx)
}

func (m *MatchStream) Match(cx *layer4.Connection) (bool, error) {
	matched, err := m.tls.Match(cx)
	if err != nil || !matched {
		return false, err
	}
	repl := cx.Context.Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	serverName, _ := repl.GetString("l4.tls.server_name")
	app := AppName(serverName, m.Domain)
	if app == "" {
		return false, nil
	}
	info, err := os.Stat(filepath.Join(m.Root, app, "common", "stream.socket"))
	if err != nil || info.Mode()&os.ModeSocket == 0 {
		return false, nil
	}
	repl.Set(AppPlaceholder, app)
	return true, nil
}

func (m *MatchStream) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next()
	if !d.Args(&m.Domain) {
		return d.ArgErr()
	}
	for d.NextBlock(0) {
		switch d.Val() {
		case "root":
			if !d.Args(&m.Root) {
				return d.ArgErr()
			}
		default:
			return d.Errf("unknown syncloud_stream option %q", d.Val())
		}
	}
	return nil
}

var (
	_ caddy.Provisioner     = (*MatchStream)(nil)
	_ layer4.ConnMatcher    = (*MatchStream)(nil)
	_ caddyfile.Unmarshaler = (*MatchStream)(nil)
)
