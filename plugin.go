package docshandler

import "github.com/curioswitch/go-docs-handler/specification"

type Plugin interface {
	GenerateSpecification() (*specification.Specification, error)
}
