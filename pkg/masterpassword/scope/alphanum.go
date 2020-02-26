package scope

func init() {
	Scopes[AlphaNumeric] = &arrayBased{
		array: []rune(alphabet + numbers),
	}
}
