package scope

const special = "`~!@#$%^&*()_+=-[]{}:;'\",./<>?\\|"

func init() {
	Scopes[KeyboardSigns] = &arrayBased{
		array: []rune(alphabet + numbers + special),
	}
}

