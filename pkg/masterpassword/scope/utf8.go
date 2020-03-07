package scope

import (
	"math/rand"
	"unicode"
)

type utf8Scope struct {
	group int
}

func (u *utf8Scope) Provide(number int) rune {
	s := rand.NewSource(int64(number))
	r := rand.New(s)

	for {
		charNo := r.Intn(u.max())
		ru := rune(charNo)
		if unicode.IsPrint(ru) {
			return ru
		}
	}
}

func (u *utf8Scope) max() int {
	values := []int{0x80, 0x800, 0x2BEF, 0x10000, 0x110000}
	return values[u.group]
}

func init() {
	Scopes[Utf8] = &parameterlessProducer{&utf8Scope{
		group: 1,
	}}
}
