package scope

const special = "`~!@#$%^&*()_+=-[]{}:;'\",./<>?\\| "

func init() {
	Scopes[KeyboardSigns] = &parameterlessProducer{&arrayBased{
		array: []rune(alphabet + numbers + special),
	}}
}

