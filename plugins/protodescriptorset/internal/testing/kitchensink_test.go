package testing

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/curioswitch/go-docs-handler/plugins/protodescriptorset"
)

//go:embed test.descriptors.pb
var testDescriptors []byte

// Generated from modified GrpcDocServiceJsonSchemaTest
//   - Mapped service to root
//   - Removed grpc-web-text mime type not supported by Connect
//   - Manually reformatted example request to follow significantly different protobuf-go
//
//go:embed armeria-spec.json
var armeriaSpecJSON string

func TestAllParameterTypesMatchesArmeria(t *testing.T) {
	p := protodescriptorset.NewPlugin(testDescriptors)
	spec, err := p.GenerateSpecification()
	if err != nil {
		t.Fatal(err)
	}
	specJSON, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, armeriaSpecJSON, string(specJSON))
}
