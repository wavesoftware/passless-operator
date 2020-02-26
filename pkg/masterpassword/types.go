package masterpassword

import scopeapi "github.com/wavesoftware/passless-operator/pkg/masterpassword/scope"

// Generator can generate a secret based on given arguments
type Generator interface {

	// Generate will create new password in predictable way
	Generate(name string, scope scopeapi.Type, counter uint, length uint8) string
}

// MasterKeyResolver resolves a master key to be used by generator
type MasterKeyResolver interface {

	// MasterKey returns a master bytes
	MasterKey() []byte
}
