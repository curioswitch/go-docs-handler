package testing

import (
	_ "embed"
	"encoding/json"
	"github.com/curioswitch/go-docs-handler/plugins/protodescriptorset"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test.descriptors.pb
var testDescriptors []byte

//go:embed armeria-spec.json
var armeriaSpecJSON string

func TestAllParameterTypes(t *testing.T) {
	t.Skip("WIP")

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
