package scope

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz"

func init() {
	Scopes[Alphabet] = &arrayBased{
		array: []rune(alphabet),
	}
}
