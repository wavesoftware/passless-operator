package scope

type parameterlessProducer struct {
		scope Scope
}

func (p *parameterlessProducer) Produce(_ string) Scope {
	return p.scope
}
