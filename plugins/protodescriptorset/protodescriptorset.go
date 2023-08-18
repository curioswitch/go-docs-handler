package protodescriptorset

import (
	"fmt"
	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/specification"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func NewPlugin(serializedDescriptors []byte) docshandler.Plugin {
	return &plugin{
		serializedDescriptors: serializedDescriptors,
	}
}

type plugin struct {
	serializedDescriptors []byte
}

func (p *plugin) GenerateSpecification() (*specification.Specification, error) {
	var descriptors descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(p.serializedDescriptors, &descriptors); err != nil {
		return nil, fmt.Errorf("unmarshaling descriptors: %w", err)
	}

	spec := &specification.Specification{}

	docstrings := make(map[string]string)

	for _, file := range descriptors.GetFile() {
		extractDocStrings(file, docstrings)
		fileDesc, _ := protodesc.NewFile(file, nil)
		for i := 0; i < fileDesc.Messages().Len(); i++ {
			spec.Structs = append(spec.Structs, convertMessage(fileDesc.Messages().Get(i), docstrings))
		}
		for i := 0; i < fileDesc.Services().Len(); i++ {
			spec.Services = append(spec.Services, convertService(fileDesc.Services().Get(i), docstrings))
		}
	}

	return spec, nil
}

func convertMessage(msg protoreflect.MessageDescriptor, docstrings map[string]string) specification.Struct {
	res := specification.Struct{
		Name: string(msg.FullName()),
	}

	if doc, ok := docstrings[string(msg.FullName())]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	for i := 0; i < msg.Fields().Len(); i++ {
		res.Fields = append(res.Fields, convertField(msg, msg.Fields().Get(i), docstrings))
	}

	return res
}

func convertField(msg protoreflect.MessageDescriptor, field protoreflect.FieldDescriptor, docstrings map[string]string) specification.Field {
	res := specification.Field{
		Name: string(field.Name()),
	}

	if msg.RequiredNumbers().Has(field.Number()) {
		res.Requirement = "required"
	} else {
		res.Requirement = "optional"
	}

	if doc, ok := docstrings[fmt.Sprintf("%s/%s", field.ContainingMessage().FullName(), field.Name())]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	var typeSignature string
	switch field.Kind() {
	case protoreflect.BoolKind:
		typeSignature = typeSignatureBool
	case protoreflect.BytesKind:
		typeSignature = typeSignatureBytes
	case protoreflect.DoubleKind:
		typeSignature = typeSignatureDouble
	case protoreflect.Fixed32Kind:
		typeSignature = typeSignatureFixed32
	case protoreflect.Fixed64Kind:
		typeSignature = typeSignatureFixed64
	case protoreflect.FloatKind:
		typeSignature = typeSignatureFloat
	case protoreflect.Int32Kind:
		typeSignature = typeSignatureInt32
	case protoreflect.Int64Kind:
		typeSignature = typeSignatureInt64
	case protoreflect.Sfixed32Kind:
		typeSignature = typeSignatureSfixed32
	case protoreflect.Sfixed64Kind:
		typeSignature = typeSignatureSfixed64
	case protoreflect.Sint32Kind:
		typeSignature = typeSignatureSint32
	case protoreflect.Sint64Kind:
		typeSignature = typeSignatureSint64
	case protoreflect.StringKind:
		typeSignature = typeSignatureString
	case protoreflect.Uint32Kind:
		typeSignature = typeSignatureUint32
	case protoreflect.Uint64Kind:
		typeSignature = typeSignatureUint64
	case protoreflect.MessageKind:
		typeSignature = string(field.Message().FullName())
	case protoreflect.GroupKind:
		// This type has been deprecated since the launch of protocol buffers to open source.
		// There is no real metadata for this in the descriptor, so we just treat as UNKNOWN
		// since it shouldn't happen in practice anyway.
		typeSignature = typeSignatureUnknown
	case protoreflect.EnumKind:
		typeSignature = string(field.Enum().FullName())
	default:
		typeSignature = typeSignatureUnknown
	}

	if field.Cardinality() == protoreflect.Repeated {
		typeSignature = fmt.Sprintf("repeated<%s>", typeSignature)
	}

	res.TypeSignature = typeSignature

	return res
}

func convertService(service protoreflect.ServiceDescriptor, docstrings map[string]string) specification.Service {
	res := specification.Service{
		Name: string(service.FullName()),
	}

	if doc, ok := docstrings[string(service.FullName())]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	for i := 0; i < service.Methods().Len(); i++ {
		res.Methods = append(res.Methods, convertMethod(service, service.Methods().Get(i), docstrings))
	}

	return res
}

func convertMethod(service protoreflect.ServiceDescriptor, method protoreflect.MethodDescriptor, docstrings map[string]string) specification.Method {
	fullName := fmt.Sprintf("%s/%s", service.FullName(), method.Name())

	endpoint := specification.Endpoint{
		DefaultMimeType:    mimeTypeGRPC,
		AvailableMimeTypes: []string{mimeTypeGRPC, mimeTypeConnectProto, mimeTypeConnectJSON},
		PathMapping:        "/" + fullName,
	}

	res := specification.Method{
		Name:                string(method.Name()),
		ID:                  fmt.Sprintf("%s/%s", fullName, "POST"),
		Endpoints:           []specification.Endpoint{endpoint},
		ReturnTypeSignature: string(method.Output().FullName()),
		Parameters: []specification.Field{
			{
				Name:          "request",
				TypeSignature: string(method.Input().FullName()),
				Requirement:   "required",
			},
		},
		HTTPMethod: "POST",
	}

	if doc, ok := docstrings[fullName]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	return res
}
