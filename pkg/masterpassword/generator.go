package masterpassword

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/wavesoftware/go-ensure"
	scopeapi "github.com/wavesoftware/passless-operator/pkg/masterpassword/scope"
	"golang.org/x/crypto/scrypt"
)

const (
	authentication = "pl.wavesoftware.masterpassword"
	identity       = "kubernetes"
)

type generator struct {
	resolver  MasterKeyResolver
	masterKey *masterKey
}

type masterKey struct {
	data    []byte
	updated time.Time
}

// NewGenerator will create new passless secret generator
func NewGenerator(resolver MasterKeyResolver) Generator {
	return &generator{
		resolver: resolver,
	}
}

func (g *generator) Generate(name, scope string, counter uint, length uint8) string {
	siteKey := calculateSiteKey(name, g.ensureMasterKey(), counter)
	return calculateSecret(siteKey, scope, length)
}

func calculateSiteKey(siteName string, masterKey []byte, counter uint) []byte {
	seed := fmt.Sprintf("%s%c%s%c",
		authentication, len(siteName), siteName, counter)
	h := hmac.New(sha256.New, []byte(seed))
	_, err := h.Write(masterKey)
	ensure.NoError(err)
	return h.Sum(nil)
}

func calculateSecret(siteKey []byte, scopeSpec string, length uint8) string {
	secret := make([]rune, 0)
	numbers := newNumberGenerator(siteKey)
	scopeType, params := processScopeSpec(scopeSpec)
	scopeProducer := scopeapi.Scopes[scopeType]
	scope := scopeProducer.Produce(params)
	for len(secret) < int(length) {
		number := numbers.next()
		r := scope.Provide(number)
		secret = append(secret, r)
	}
	return string(secret)
}

func processScopeSpec(scopeSpec string) (scopeapi.Type, string) {
	parts := strings.SplitN(scopeSpec, ":", 2)
	t := scopeapi.Type(parts[0])
	if len(parts) == 2 {
		return t, parts[1]
	}
	return t, ""
}

func (g *generator) ensureMasterKey() []byte {
	now := time.Now()
	lastTime := now.Add(-1 * time.Minute * 15)
	if g.masterKey == nil || g.masterKey.updated.Before(lastTime) {
		data := g.resolver.MasterKey()
		g.masterKey = &masterKey{
			data:    calculateMasterKey(data),
			updated: now,
		}
	}
	return g.masterKey.data
}

func calculateMasterKey(secret []byte) []byte {
	key := secret
	salt := authentication + string(rune(len(identity))) + identity
	cost := 32_768
	blocksize := 8
	parallelization := 2
	length := 64
	return crypt(key, salt, cost, blocksize, parallelization, length)
}

func crypt(key []byte, salt string, cost, blocksize, parallelization, length int) []byte {
	dk, err := scrypt.Key(key, []byte(salt), cost, blocksize, parallelization, length)
	ensure.NoError(err)
	return dk
}
