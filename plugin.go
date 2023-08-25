package docshandler

import "github.com/curioswitch/go-docs-handler/specification"

// Plugin generates a specifications for services. Different frameworks
// used for service definition will each have a plugin for introspecting
// registered services and generating documentation.
type Plugin interface {
	// GenerateSpecification generates a specification for the services.
	GenerateSpecification() (*specification.Specification, error)
}
