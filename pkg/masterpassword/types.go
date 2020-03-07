package masterpassword

type ScopeType string

const (
	Numeric       ScopeType = "num"
	Alphabet      ScopeType = "alpha"
	AlphaNumeric  ScopeType = "alnum"
	EasyForHuman  ScopeType = "human"
	KeyboardSigns ScopeType = "keys"
	Utf8          ScopeType = "utf8"
	Listing       ScopeType = "list"
)

type PassLessGenerator interface {
	Generate(identity, name string, scope ScopeType, counter uint, length uint8) []byte
}
