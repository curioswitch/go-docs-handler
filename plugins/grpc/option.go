package grpc

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

type Option func(*config)

func WithSerializedDescriptors(serializedDescriptors []byte) Option {
	return func(c *config) {
		c.serializedDescriptors = serializedDescriptors
	}
}

func WithExampleRequests(method string, request proto.Message, more ...proto.Message) Option {
	return func(c *config) {
		c.exampleRequests[method] = append([]proto.Message{request}, more...)
	}
}
