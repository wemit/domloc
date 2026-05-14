package caddy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wemit/domloc/internal/registry"
)

func TestRenderCaddyfile_HTTPS(t *testing.T) {
	routes := []registry.Route{{Domain: "app.test", Port: 3000, HTTPS: true}}
	out, err := renderCaddyfile(routes, "/tmp/caddy-data")
	require.NoError(t, err)
	assert.Contains(t, out, "app.test {")
	assert.Contains(t, out, "reverse_proxy localhost:3000")
	assert.Contains(t, out, "tls internal")
	assert.Contains(t, out, "root /tmp/caddy-data")
}

func TestRenderCaddyfile_HTTP(t *testing.T) {
	routes := []registry.Route{{Domain: "api.test", Port: 4000, HTTPS: false}}
	out, err := renderCaddyfile(routes, "/tmp/caddy-data")
	require.NoError(t, err)
	assert.Contains(t, out, "http://api.test {")
	assert.NotContains(t, out, "tls internal")
}

func TestRenderCaddyfile_Mixed(t *testing.T) {
	routes := []registry.Route{
		{Domain: "app.test", Port: 3000, HTTPS: true},
		{Domain: "api.test", Port: 4000, HTTPS: false},
	}
	out, err := renderCaddyfile(routes, "/tmp/caddy-data")
	require.NoError(t, err)
	assert.Contains(t, out, "app.test {")
	assert.Contains(t, out, "http://api.test {")
	assert.Contains(t, out, "reverse_proxy localhost:3000")
	assert.Contains(t, out, "reverse_proxy localhost:4000")
}

func TestRenderCaddyfile_Empty(t *testing.T) {
	out, err := renderCaddyfile(nil, "/tmp/caddy-data")
	require.NoError(t, err)
	assert.Contains(t, out, "storage file_system")
	assert.NotContains(t, out, "reverse_proxy")
}
