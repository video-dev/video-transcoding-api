package provider

// Descriptor fully describes a provider.
//
// It contains the name of the provider, along with its current heath status
// and its capabilities.
type Descriptor struct {
	Name         string       `json:"name"`
	Capabilities Capabilities `json:"capabilities"`
	Health       Health       `json:"health"`
}

// Capabilities describes the available features in the provider. It specificie
// which input and output formats the provider supports, along with
// supported destinations.
type Capabilities struct {
	InputFormats  []string `json:"input"`
	OutputFormats []string `json:"output"`
	Destinations  []string `json:"destinations"`
}

// Health describes the current health status of the provider. If indicates
// whether the provider is healthy or not, and if it's not healthy, it includes
// a message explaining what's wrong.
type Health struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}
