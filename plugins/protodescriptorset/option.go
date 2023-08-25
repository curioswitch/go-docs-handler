package protodescriptorset

import "google.golang.org/protobuf/proto"

type config struct {
	exampleRequests map[string][]proto.Message
}

func newConfig() *config {
	return &config{
		exampleRequests: make(map[string][]proto.Message),
	}
}

type Option func(*config)

func WithExampleRequests(method string, requests []proto.Message) Option {
	return func(c *config) {
		c.exampleRequests[method] = requests
	}
}
