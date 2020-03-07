package scope

func init() {
	Scopes[AlphaNumeric] = &parameterlessProducer{&arrayBased{
		array: []rune(alphabet + numbers),
	}}
}
