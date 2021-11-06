package metadata_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesoftware/passless-operator/pkg/metadata"
)

func TestVersion(t *testing.T) {
	assert.NotEmpty(t, metadata.Version)
	assert.NotEqual(t, "v0.0.0", metadata.Version)
}
