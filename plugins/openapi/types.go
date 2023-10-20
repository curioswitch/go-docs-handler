package openapi

import (
	"encoding/json"

	"github.com/curioswitch/go-docs-handler/specification"
)

type jsonTypeSignature int

const (
	jsonTypeSignatureBoolean jsonTypeSignature = iota
	jsonTypeSignatureInt
	jsonTypeSignatureUint
	jsonTypeSignatureInt32
	jsonTypeSignatureInt64
	jsonTypeSignatureUint32
	jsonTypeSignatureUint64
	jsonTypeSignatureNumber
	jsonTypeSignatureFloat
	jsonTypeSignatureDouble
	jsonTypeSignatureString
	jsonTypeSignatureUnknown
)

func (j jsonTypeSignature) Type() specification.TypeSignatureType {
	return specification.TypeSignatureTypeBase
}

func (j jsonTypeSignature) Signature() string {
	switch j {
	case jsonTypeSignatureBoolean:
		return "boolean"
	case jsonTypeSignatureInt:
		return "int"
	case jsonTypeSignatureUint:
		return "uint"
	case jsonTypeSignatureInt32:
		return "int32"
	case jsonTypeSignatureInt64:
		return "int64"
	case jsonTypeSignatureUint32:
		return "uint32"
	case jsonTypeSignatureUint64:
		return "uint64"
	case jsonTypeSignatureNumber:
		return "number"
	case jsonTypeSignatureFloat:
		return "float"
	case jsonTypeSignatureDouble:
		return "double"
	case jsonTypeSignatureString:
		return "string"
	case jsonTypeSignatureUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

func (j jsonTypeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Signature())
}
