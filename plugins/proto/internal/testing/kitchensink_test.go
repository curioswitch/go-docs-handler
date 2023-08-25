package testing

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/kinbiko/jsonassert"

	protodocs "github.com/curioswitch/go-docs-handler/plugins/proto"
)

//go:embed test.descriptors.pb
var testDescriptors []byte

// Generated from modified GrpcDocServiceJsonSchemaTest
//   - Mapped service to root
//   - Removed grpc-web-text mime type not supported by Connect
//   - Example request replaced with presence check because JSON strings in Go are not stable
//
//go:embed armeria-spec.json
var armeriaSpecJSON string

func TestAllParameterTypesMatchesArmeria(t *testing.T) {
	p := protodocs.NewPlugin("armeria.grpc.testing.TestService",
		protodocs.WithSerializedDescriptors(testDescriptors))
	spec, err := p.GenerateSpecification()
	if err != nil {
		t.Fatal(err)
	}
	specJSON, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(specJSON), armeriaSpecJSON)
}
