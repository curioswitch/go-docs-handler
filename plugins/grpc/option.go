package grpcdocs

import "google.golang.org/protobuf/proto"

type config struct {
	serializedDescriptors []byte
	exampleRequests       map[string][]proto.Message
}

func newConfig() *config {
	return &config{
		exampleRequests: make(map[string][]proto.Message),
	}
}

// Option is a configuration option for NewPlugin.
type Option func(*config)

// WithSerializedDescriptors registers serialized proto descriptors to use for extracting
// docstrings in documentation.
func WithSerializedDescriptors(serializedDescriptors []byte) Option {
	return func(c *config) {
		c.serializedDescriptors = serializedDescriptors
	}
}

// WithExampleRequests registers example requests for a method which can be selected in the
// debug form of the docs page. The method name should be in the format "service/method", e.g.,
// "greet.GreetService/Greet".
func WithExampleRequests(method string, request proto.Message, more ...proto.Message) Option {
	return func(c *config) {
		c.exampleRequests[method] = append([]proto.Message{request}, more...)
	}
}
