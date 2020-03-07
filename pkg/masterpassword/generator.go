package masterpassword

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(
	identity, name string, scope ScopeType, counter uint, length uint8,
	) []byte {

	return []byte{45, 65}
}
