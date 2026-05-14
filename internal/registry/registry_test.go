package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAndFind(t *testing.T) {
	r := &Registry{}
	r.Add(Route{Domain: "app.test", Port: 3000})
	got, ok := r.Find("app.test")
	require.True(t, ok)
	assert.Equal(t, 3000, got.Port)
}

func TestAddOverwrites(t *testing.T) {
	r := &Registry{}
	r.Add(Route{Domain: "app.test", Port: 3000})
	r.Add(Route{Domain: "app.test", Port: 4000})
	assert.Len(t, r.Routes, 1)
	got, _ := r.Find("app.test")
	assert.Equal(t, 4000, got.Port)
}

func TestRemove(t *testing.T) {
	r := &Registry{}
	r.Add(Route{Domain: "app.test", Port: 3000})
	assert.True(t, r.Remove("app.test"))
	assert.False(t, r.Remove("app.test"))
	assert.Empty(t, r.Routes)
}

func TestIsHTTPSDefault(t *testing.T) {
	r := &Registry{}
	assert.True(t, r.IsHTTPSDefault(), "nil pointer = true (unset = HTTPS)")

	f := false
	r.HTTPSDefault = &f
	assert.False(t, r.IsHTTPSDefault())

	tr := true
	r.HTTPSDefault = &tr
	assert.True(t, r.IsHTTPSDefault())
}

func TestSaveLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "routes.json")

	r := &Registry{path: path, Routes: []Route{
		{Domain: "app.test", Port: 3000, HTTPS: true},
	}}
	require.NoError(t, r.Save())

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var r2 Registry
	require.NoError(t, json.Unmarshal(data, &r2))
	require.Len(t, r2.Routes, 1)
	assert.Equal(t, "app.test", r2.Routes[0].Domain)
}
