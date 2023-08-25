package protodocs

import (
	"encoding/json"

	"github.com/curioswitch/go-docs-handler/specification"
)

type protoTypeSignature int

const (
	protoTypeSignatureBool protoTypeSignature = iota
	protoTypeSignatureInt32
	protoTypeSignatureInt64
	protoTypeSignatureUint32
	protoTypeSignatureUint64
	protoTypeSignatureSint32
	protoTypeSignatureSint64
	protoTypeSignatureFixed32
	protoTypeSignatureFixed64
	protoTypeSignatureSfixed32
	protoTypeSignatureSfixed64
	protoTypeSignatureFloat
	protoTypeSignatureDouble
	protoTypeSignatureString
	protoTypeSignatureBytes
	protoTypeSignatureUnknown
)

func (p protoTypeSignature) Type() specification.TypeSignatureType {
	return specification.TypeSignatureTypeBase
}

func (p protoTypeSignature) Signature() string {
	switch p {
	case protoTypeSignatureBool:
		return "bool"
	case protoTypeSignatureInt32:
		return "int32"
	case protoTypeSignatureInt64:
		return "int64"
	case protoTypeSignatureUint32:
		return "uint32"
	case protoTypeSignatureUint64:
		return "uint64"
	case protoTypeSignatureSint32:
		return "sint32"
	case protoTypeSignatureSint64:
		return "sint64"
	case protoTypeSignatureFixed32:
		return "fixed32"
	case protoTypeSignatureFixed64:
		return "fixed64"
	case protoTypeSignatureSfixed32:
		return "sfixed32"
	case protoTypeSignatureSfixed64:
		return "sfixed64"
	case protoTypeSignatureFloat:
		return "float"
	case protoTypeSignatureDouble:
		return "double"
	case protoTypeSignatureString:
		return "string"
	case protoTypeSignatureBytes:
		return "bytes"
	default:
		return "unknown"
	}
}

func (p protoTypeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Signature())
}
