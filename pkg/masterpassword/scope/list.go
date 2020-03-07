package scope

type listProducer struct{}

func (p *listProducer) Produce(params string) Scope {
	return &arrayBased{
		array: []rune(params),
	}
}

func init() {
	Scopes[Listing] = &listProducer{}
}
