package scope

func init() {
	Scopes[EasyForHuman] = &parameterlessProducer{&arrayBased{
		// Ref: https://stackoverflow.com/a/55634/844449
		array: []rune("!#%+23456789:=?@ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"),
	}}
}
