package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderDnsmasqConf_Single(t *testing.T) {
	out, err := renderDnsmasqConf([]string{"test"})
	require.NoError(t, err)
	assert.Contains(t, out, "port=5300")
	assert.Contains(t, out, "listen-address=127.0.0.1")
	assert.Contains(t, out, "address=/.test/127.0.0.1")
}

func TestRenderDnsmasqConf_Multiple(t *testing.T) {
	out, err := renderDnsmasqConf([]string{"test", "loc"})
	require.NoError(t, err)
	assert.Contains(t, out, "address=/.test/127.0.0.1")
	assert.Contains(t, out, "address=/.loc/127.0.0.1")
}

func TestRenderDnsmasqConf_Empty(t *testing.T) {
	out, err := renderDnsmasqConf(nil)
	require.NoError(t, err)
	assert.Contains(t, out, "port=5300")
	assert.NotContains(t, out, "address=/.")
}
