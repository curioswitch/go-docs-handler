package protodescriptorset

import (
	"fmt"
	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/specification"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"net/http"
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
			structs, enums := convertMessage(fileDesc.Messages().Get(i), docstrings)
			spec.Structs = append(spec.Structs, structs...)
			spec.Enums = append(spec.Enums, enums...)
		}
		for i := 0; i < fileDesc.Enums().Len(); i++ {
			spec.Enums = append(spec.Enums, convertEnum(fileDesc.Enums().Get(i), docstrings))
		}
		for i := 0; i < fileDesc.Services().Len(); i++ {
			spec.Services = append(spec.Services, convertService(fileDesc.Services().Get(i), docstrings))
		}
	}

	return spec, nil
}

func convertMessage(msg protoreflect.MessageDescriptor, docstrings map[string]string) ([]specification.Struct, []specification.Enum) {
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

	var enums []specification.Enum
	structs := []specification.Struct{res}

	for i := 0; i < msg.Enums().Len(); i++ {
		enums = append(enums, convertEnum(msg.Enums().Get(i), docstrings))
	}

	for i := 0; i < msg.Messages().Len(); i++ {
		s, e := convertMessage(msg.Messages().Get(i), docstrings)
		structs = append(structs, s...)
		enums = append(enums, e...)
	}

	return structs, enums
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

	res.TypeSignature = fieldTypeSignature(field)

	return res
}

func fieldTypeSignature(field protoreflect.FieldDescriptor) specification.TypeSignature {
	if field.IsMap() {
		return specification.NewMapTypeSignature(
			fieldTypeSignature(field.MapKey()),
			fieldTypeSignature(field.MapValue()),
		)
	}

	var typeSignature specification.TypeSignature
	switch field.Kind() {
	case protoreflect.BoolKind:
		typeSignature = protoTypeSignatureBool
	case protoreflect.BytesKind:
		typeSignature = protoTypeSignatureBytes
	case protoreflect.DoubleKind:
		typeSignature = protoTypeSignatureDouble
	case protoreflect.Fixed32Kind:
		typeSignature = protoTypeSignatureFixed32
	case protoreflect.Fixed64Kind:
		typeSignature = protoTypeSignatureFixed64
	case protoreflect.FloatKind:
		typeSignature = protoTypeSignatureFloat
	case protoreflect.Int32Kind:
		typeSignature = protoTypeSignatureInt32
	case protoreflect.Int64Kind:
		typeSignature = protoTypeSignatureInt64
	case protoreflect.Sfixed32Kind:
		typeSignature = protoTypeSignatureSfixed32
	case protoreflect.Sfixed64Kind:
		typeSignature = protoTypeSignatureSfixed64
	case protoreflect.Sint32Kind:
		typeSignature = protoTypeSignatureSint32
	case protoreflect.Sint64Kind:
		typeSignature = protoTypeSignatureSint64
	case protoreflect.StringKind:
		typeSignature = protoTypeSignatureString
	case protoreflect.Uint32Kind:
		typeSignature = protoTypeSignatureUint32
	case protoreflect.Uint64Kind:
		typeSignature = protoTypeSignatureUint64
	case protoreflect.MessageKind:
		typeSignature = specification.NewStructTypeSignature(string(field.Message().FullName()))
	case protoreflect.GroupKind:
		// This type has been deprecated since the launch of protocol buffers to open source.
		// There is no real metadata for this in the descriptor, so we just treat as UNKNOWN
		// since it shouldn't happen in practice anyway.
		typeSignature = protoTypeSignatureUnknown
	case protoreflect.EnumKind:
		typeSignature = specification.NewEnumTypeSignature(string(field.Enum().FullName()))
	default:
		typeSignature = protoTypeSignatureUnknown
	}

	if field.Cardinality() == protoreflect.Repeated {
		typeSignature = specification.NewIterableTypeSignature("repeated", typeSignature)
	}

	return typeSignature
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
		ReturnTypeSignature: specification.NewStructTypeSignature(string(method.Output().FullName())),
		Parameters: []specification.Field{
			{
				Name:          "request",
				TypeSignature: specification.NewStructTypeSignature(string(method.Input().FullName())),
				Requirement:   "required",
			},
		},
		UseParameterAsRoot: true,
		HTTPMethod:         http.MethodPost,
	}

	if doc, ok := docstrings[fullName]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	return res
}

func convertEnum(enum protoreflect.EnumDescriptor, docstrings map[string]string) specification.Enum {
	res := specification.Enum{
		Name: string(enum.FullName()),
	}

	if doc, ok := docstrings[string(enum.FullName())]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	for i := 0; i < enum.Values().Len(); i++ {
		res.Values = append(res.Values, convertEnumValue(enum, enum.Values().Get(i), docstrings))
	}

	return res
}

func convertEnumValue(enum protoreflect.EnumDescriptor, value protoreflect.EnumValueDescriptor, docstrings map[string]string) specification.Value {
	intVal := int(value.Number())
	res := specification.Value{
		Name:     string(value.Name()),
		IntValue: &intVal,
	}

	if doc, ok := docstrings[fmt.Sprintf("%s/%s", enum.FullName(), value.Name())]; ok {
		res.DescriptionInfo.DocString = doc
		res.DescriptionInfo.Markup = "NONE"
	}

	return res
}
