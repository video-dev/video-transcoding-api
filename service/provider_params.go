package service

// swagger:parameters getProvider deleteProvider
type getProviderInput struct {
	// in: path
	// required: true
	Name string `json:"name"`
}

func (p *getProviderInput) loadParams(paramsMap map[string]string) {
	p.Name = paramsMap["name"]
}
