package scope

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz"

func init() {
	Scopes[Alphabet] = &parameterlessProducer{&arrayBased{
		array: []rune(alphabet),
	}}
}
