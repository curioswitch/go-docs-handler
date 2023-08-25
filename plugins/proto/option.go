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

// Option is a configuration option for NewPlugin.
type Option func(*config)

// WithAdditionalService registers an additional service to generate documentation for.
func WithAdditionalService(service string) Option {
	return func(c *config) {
		c.services = append(c.services, service)
	}
}

// WithSerializedDescriptors registers serialized proto descriptors to use for extracting
// docstrings in documentation.
func WithSerializedDescriptors(serializedDescriptors []byte) Option {
	return func(c *config) {
		c.serializedDescriptors = serializedDescriptors
	}
}

// WithExampleRequests registers example requests for a method which can be selected in the
// debug form of the docs page. The method name should be in the format "service/method", e.g.,
// "greet.GreetService/Greet". Connect users can use the constants for the procedures, e.g.,
// greetconnect.GreetServiceGreetProcedure.
func WithExampleRequests(method string, request proto.Message, more ...proto.Message) Option {
	// Support connect procedure format.
	if method[0] == '/' {
		method = method[1:]
	}

	return func(c *config) {
		c.exampleRequests[method] = append([]proto.Message{request}, more...)
	}
}
