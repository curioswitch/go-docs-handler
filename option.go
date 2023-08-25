package docshandler

type config struct {
	plugins                 []Plugin
	injectedScriptSuppliers []func() string
}

// Option is a configuration option for the docs handler.
type Option func(*config)

func newConfig(plugin Plugin) *config {
	return &config{
		plugins: []Plugin{plugin},
	}
}

// WithAdditionalPlugin registers an additional plugin to serve documentation for.
func WithAdditionalPlugin(plugin Plugin) Option {
	return func(c *config) {
		c.plugins = append(c.plugins, plugin)
	}
}

// WithInjectedScriptSupplier registers a supplier for a Javascript that is injected
// into the docs page on load. This can be used to add an auth token to debug requests.
func WithInjectedScriptSupplier(supplier func() string) Option {
	return func(c *config) {
		c.injectedScriptSuppliers = append(c.injectedScriptSuppliers, supplier)
	}
}
