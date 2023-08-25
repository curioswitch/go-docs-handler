package protodocs

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/specification"
)

func NewPlugin(service string, opts ...Option) docshandler.Plugin {
	c := newConfig(service)
	for _, opt := range opts {
		opt(c)
	}
	return &plugin{
		config: c,
	}
}

type plugin struct {
	config *config
}

func (p *plugin) GenerateSpecification() (*specification.Specification, error) {
	spec := &specification.Specification{}

	var services []protoreflect.ServiceDescriptor
	for _, service := range p.config.services {
		desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(service))
		if err != nil {
			return nil, fmt.Errorf("proto: unable to find descriptor for service '%s': %w", service, err)
		}
		if svc, ok := desc.(protoreflect.ServiceDescriptor); ok {
			services = append(services, svc)
		} else {
			return nil, fmt.Errorf("proto: descriptor for service '%s' is not a service", service)
		}
	}

	docstrings := make(map[string]string)
	if p.config.serializedDescriptors != nil {
		var descriptors descriptorpb.FileDescriptorSet
		if err := proto.Unmarshal(p.config.serializedDescriptors, &descriptors); err != nil {
			return nil, fmt.Errorf("unmarshaling descriptors: %w", err)
		}
		for _, file := range descriptors.GetFile() {
			extractDocStrings(file, docstrings)
		}
	}

	msgs := make(map[string]protoreflect.MessageDescriptor)
	enums := make(map[string]protoreflect.EnumDescriptor)

	// First find what we need to generate for.
	for _, svc := range services {
		findServiceMessagesAndEnums(svc, msgs, enums)
	}

	// Then generate.
	for _, msg := range msgs {
		spec.Structs = append(spec.Structs, convertMessage(msg, docstrings))
	}
	sort.SliceStable(spec.Structs, func(i, j int) bool {
		return spec.Structs[i].Name < spec.Structs[j].Name
	})

	for _, enum := range enums {
		spec.Enums = append(spec.Enums, convertEnum(enum, docstrings))
	}
	sort.SliceStable(spec.Enums, func(i, j int) bool {
		return spec.Enums[i].Name < spec.Enums[j].Name
	})

	for _, svc := range services {
		s, err := p.convertService(svc, docstrings)
		if err != nil {
			return nil, err
		}
		spec.Services = append(spec.Services, s)
	}
	sort.SliceStable(spec.Services, func(i, j int) bool {
		return spec.Services[i].Name < spec.Services[j].Name
	})

	return spec, nil
}

func findServiceMessagesAndEnums(svc protoreflect.ServiceDescriptor, msgs map[string]protoreflect.MessageDescriptor, enums map[string]protoreflect.EnumDescriptor) {
	for i := 0; i < svc.Methods().Len(); i++ {
		m := svc.Methods().Get(i)
		findMessageDependencies(m.Input(), msgs, enums)
		findMessageDependencies(m.Output(), msgs, enums)
	}
}

func findMessageDependencies(msg protoreflect.MessageDescriptor, msgs map[string]protoreflect.MessageDescriptor, enums map[string]protoreflect.EnumDescriptor) {
	if _, ok := msgs[string(msg.FullName())]; ok || msg.IsMapEntry() {
		return
	}
	msgs[string(msg.FullName())] = msg
	for i := 0; i < msg.Fields().Len(); i++ {
		field := msg.Fields().Get(i)
		if field.Message() != nil {
			findMessageDependencies(field.Message(), msgs, enums)
		} else if field.Enum() != nil {
			enumName := string(field.Enum().FullName())
			if _, ok := enums[enumName]; !ok {
				enums[enumName] = field.Enum()
			}
		}
	}

	for i := 0; i < msg.Enums().Len(); i++ {
		enum := msg.Enums().Get(i)
		enumName := string(enum.FullName())
		if _, ok := enums[enumName]; !ok {
			enums[enumName] = enum
		}
	}

	for i := 0; i < msg.Messages().Len(); i++ {
		m := msg.Messages().Get(i)
		findMessageDependencies(m, msgs, enums)
	}
}

func convertMessage(msg protoreflect.MessageDescriptor, docstrings map[string]string) specification.Struct {
	res := specification.Struct{
		Name: string(msg.FullName()),
	}

	if doc, ok := docstrings[string(msg.FullName())]; ok {
		res.DescriptionInfo.DocString = doc
	}
	res.DescriptionInfo.Markup = "NONE"

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
		res.Requirement = "REQUIRED"
	} else {
		res.Requirement = "OPTIONAL"
	}

	if doc, ok := docstrings[fmt.Sprintf("%s/%s", field.ContainingMessage().FullName(), field.Name())]; ok {
		res.DescriptionInfo.DocString = doc
	}
	res.DescriptionInfo.Markup = "NONE"

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

func (p *plugin) convertService(service protoreflect.ServiceDescriptor, docstrings map[string]string) (specification.Service, error) {
	res := specification.Service{
		Name: string(service.FullName()),
	}

	if doc, ok := docstrings[string(service.FullName())]; ok {
		res.DescriptionInfo.DocString = doc
	}
	res.DescriptionInfo.Markup = "NONE"

	for i := 0; i < service.Methods().Len(); i++ {
		m, err := p.convertMethod(service, service.Methods().Get(i), docstrings)
		if err != nil {
			return specification.Service{}, err
		}
		res.Methods = append(res.Methods, m)
	}

	return res, nil
}

func (p *plugin) convertMethod(service protoreflect.ServiceDescriptor, method protoreflect.MethodDescriptor, docstrings map[string]string) (specification.Method, error) {
	fullName := fmt.Sprintf("%s/%s", service.FullName(), method.Name())

	endpoint := specification.Endpoint{
		// Add mime types supported by Connect which should be most usage of this plugin.
		AvailableMimeTypes: []string{mimeTypeGRPCJSON, mimeTypeGRPCProto, mimeTypeGRPCWebJSON, mimeTypeGRPCWebProto, mimeTypeConnectJSON, mimeTypeConnectProto},
		HostnamePattern:    "*",
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
				Requirement:   "REQUIRED",
				DescriptionInfo: specification.DescriptionInfo{
					Markup: "NONE",
				},
			},
		},
		UseParameterAsRoot: true,
		HTTPMethod:         http.MethodPost,
	}

	pjson := protojson.MarshalOptions{Multiline: true, EmitUnpopulated: true, UseProtoNames: true}

	if reqs, ok := p.config.exampleRequests[fullName]; ok {
		for _, req := range reqs {
			reqJSON, err := pjson.Marshal(req)
			if err != nil {
				return specification.Method{}, fmt.Errorf("protodescriptorset: marshaling example request: %w", err)
			}
			res.ExampleRequests = append(res.ExampleRequests, string(reqJSON))
		}
	}

	prototypeJSON, err := pjson.Marshal(dynamicpb.NewMessage(method.Input()))
	if err != nil {
		return specification.Method{}, fmt.Errorf("protodescriptorset: marshaling prototype request: %w", err)
	}
	if !bytes.Equal(prototypeJSON, []byte("{\n}")) && !bytes.Equal(prototypeJSON, []byte("{}")) {
		res.ExampleRequests = append(res.ExampleRequests, string(prototypeJSON))
	}

	if doc, ok := docstrings[fullName]; ok {
		res.DescriptionInfo.DocString = doc
	}
	res.DescriptionInfo.Markup = "NONE"

	return res, nil
}

func convertEnum(enum protoreflect.EnumDescriptor, docstrings map[string]string) specification.Enum {
	res := specification.Enum{
		Name: string(enum.FullName()),
	}

	if doc, ok := docstrings[string(enum.FullName())]; ok {
		res.DescriptionInfo.DocString = doc
	}
	res.DescriptionInfo.Markup = "NONE"

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
	}
	res.DescriptionInfo.Markup = "NONE"

	return res
}
