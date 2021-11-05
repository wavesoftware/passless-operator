package passless_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesoftware/passless-operator/pkg/passless"
	"github.com/wavesoftware/passless-operator/test"
)

func TestResolverMasterKey(t *testing.T) {
	// when
	client, err := test.NewClient()
	if err != nil {
		t.Skip(err)
	}
	resolver := passless.Resolver{Client: client, Context: context.TODO()}

	// when
	masterKey := resolver.MasterKey()
	masterKey2 := resolver.MasterKey()

	// then
	assert.NotEmpty(t, masterKey)
	assert.Greater(t, len(masterKey), 10)
	assert.Equal(t, masterKey, masterKey2)
}
