package masterpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	scopeapi "github.com/wavesoftware/passless-operator/pkg/masterpassword/scope"
)

type mockResolver struct {
	val string
}

func (t *mockResolver) MasterKey() []byte {
	return []byte(t.val)
}

func TestGenerator(t *testing.T) {
	resolver := &mockResolver{val: "mock contents"}
	g := NewGenerator(resolver)

	secret := g.Generate("db-secret", scopeapi.Utf8, 1, 10)

	assert.Equal(t, "˸ҰݴާފӇͰƍߕѫ", secret)
}
