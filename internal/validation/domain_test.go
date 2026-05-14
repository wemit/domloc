package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomain(t *testing.T) {
	cases := []struct {
		domain string
		valid  bool
	}{
		{"app.test", true},
		{"*.foo.test", true},
		{"my-app.local", true},
		{"", false},
		{"notadomain", false},
		{"has space.test", false},
	}
	for _, tc := range cases {
		err := Domain(tc.domain)
		if tc.valid {
			assert.NoError(t, err, tc.domain)
		} else {
			assert.Error(t, err, tc.domain)
		}
	}
}

func TestPort(t *testing.T) {
	assert.NoError(t, Port(3000))
	assert.NoError(t, Port(1))
	assert.NoError(t, Port(65535))
	assert.Error(t, Port(0))
	assert.Error(t, Port(65536))
}

func TestIsWildcard(t *testing.T) {
	assert.True(t, IsWildcard("*.foo.test"))
	assert.False(t, IsWildcard("foo.test"))
}

func TestTLD(t *testing.T) {
	assert.Equal(t, "test", TLD("app.test"))
	assert.Equal(t, "test", TLD("*.foo.test"))
}
