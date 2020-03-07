package masterpassword

import (
	"fmt"
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

	tests := []struct {
		name, scope string
		counter     uint
		length      uint8
		expected    string
	}{
		{"db-secret", string(scopeapi.AlphaNumeric), 1, 10, "PzztNnez2k"},
		{"db-sec", string(scopeapi.AlphaNumeric), 1, 10, "K9w0z0gnHB"},
		{"db-sec", string(scopeapi.AlphaNumeric), 1, 14, "K9w0z0gnHB4n8Z"},
		{"db-sec", string(scopeapi.AlphaNumeric), 2, 10, "iW5Gqtvilo"},

		{"secret", string(scopeapi.Utf8), 1, 10, "ĥ&ѤͶЍӲޱވרƂ"},
		{"secret", string(scopeapi.Alphabet), 1, 10, "UNpDmPETYj"},
		{"secret", string(scopeapi.AlphaNumeric), 1, 10, "gX2b9xeD96"},
		{"secret", string(scopeapi.Numeric), 1, 10, "5264787036"},
		{"secret", string(scopeapi.EasyForHuman), 1, 10, "E#tvGSwn6z"},
		{"secret", string(scopeapi.KeyboardSigns), 1, 10, "_09;f.u!:y"},
		{"secret", "list:abc", 1, 10, "caaabbaabb"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("for name=%s scope=%s counter=%d length=%d",
			tt.name, tt.scope, tt.counter, tt.length)
		t.Run(testName, func(t *testing.T) {
			secret := g.Generate(tt.name, tt.scope, tt.counter, tt.length)

			assert.Equal(t, tt.expected, secret)
		})
	}
}
