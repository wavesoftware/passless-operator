package scope

// Type is scope of passless secret, it's type
type Type string

const (
	// Numeric type contains only numbers
	Numeric Type = "num"

	// Alphabet type contains lower and upper chars from US ascii table
	Alphabet Type = "alpha"

	// AlphaNumeric contains lower and upper chars from US ascii table and numbers
	AlphaNumeric Type = "alnum"

	// EasyForHuman contains only a subset of chars and numbers that are easy to distinguish
	EasyForHuman Type = "human"

	// KeyboardSigns contains all sings from regular US keyboard: chars, numbers, and special signs
	KeyboardSigns Type = "keys"

	// Utf8 contains almost all printable UTF-8 characters
	Utf8 Type = "utf8"

	// Listing will create secrets from user provided symbols
	Listing Type = "list"
)

// Producer creates a scope for given params
type Producer interface {
	// Produce a scope for given params
	Produce(params string) Scope
}

// Scope is used to generate secret values
type Scope interface {
	// Provide a rune by number given
	Provide(number int) rune
}

// Scopes holds all implementations of scope type
var Scopes = make(map[Type]Producer)
