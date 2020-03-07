package scope

const numbers = "1234567890"

func init() {
	Scopes[Numeric] = &parameterlessProducer{&arrayBased{
		array: []rune(numbers),
	}}
}
