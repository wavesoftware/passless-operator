package scope

type arrayBased struct {
	array []rune
}

func (a *arrayBased) Provide(num int) rune {
	idx := num % len(a.array)
	return a.array[idx]
}
