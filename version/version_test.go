package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.NotEmpty(t, Version)
	assert.NotEqual(t, "v0.0.0", Version)
}
