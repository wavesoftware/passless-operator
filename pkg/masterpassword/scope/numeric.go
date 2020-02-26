package scope

const numbers = "1234567890"

func init() {
	Scopes[Numeric] = &arrayBased{
		array: []rune(numbers),
	}
}
