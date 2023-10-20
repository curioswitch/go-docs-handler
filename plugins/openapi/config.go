package openapi

type config struct {
	spec            []byte
	exampleRequests map[string][]any
}

// Option is a configuration option for NewPlugin.
type Option func(*config)

func newConfig(spec []byte) *config {
	return &config{
		spec:            spec,
		exampleRequests: make(map[string][]any),
	}
}

// WithExampleRequests registers example requests for a method which can be selected in the
// debug form of the docs page.
func WithExampleRequests(operationID string, req any, more ...any) Option {
	requests := append([]any{req}, more...)
	return func(c *config) {
		c.exampleRequests[operationID] = requests
	}
}
