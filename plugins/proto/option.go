package protodocs

import "google.golang.org/protobuf/proto"

type config struct {
	services              []string
	serializedDescriptors []byte
	exampleRequests       map[string][]proto.Message
}

func newConfig(service string) *config {
	return &config{
		services:        []string{service},
		exampleRequests: make(map[string][]proto.Message),
	}
}

type Option func(*config)

func WithAdditionalService(service string) Option {
	return func(c *config) {
		c.services = append(c.services, service)
	}
}

func WithSerializedDescriptors(serializedDescriptors []byte) Option {
	return func(c *config) {
		c.serializedDescriptors = serializedDescriptors
	}
}

func WithExampleRequests(method string, request proto.Message, more ...proto.Message) Option {
	// Support connect procedure format.
	if method[0] == '/' {
		method = method[1:]
	}

	return func(c *config) {
		c.exampleRequests[method] = append([]proto.Message{request}, more...)
	}
}
