package protodescriptorset

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

func extractDocStrings(file *descriptorpb.FileDescriptorProto, docstrings map[string]string) {
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		if len(loc.GetLeadingComments()) == 0 {
			continue
		}
		name, ok := getFullName(file, loc.GetPath())
		if !ok {
			continue
		}
		docstrings[name] = loc.GetLeadingComments()
	}
}

func getFullName(descriptor *descriptorpb.FileDescriptorProto, path []int32) (string, bool) {
	name := strings.Builder{}
	name.WriteString(descriptor.GetPackage())

	switch path[0] {
	case fileDescriptorProtoMessageTypeFieldNumber:
		msg := descriptor.GetMessageType()[path[1]]
		if !appendMessageToFullName(msg, path, &name) {
			return "", false
		}
	case fileDescriptorProtoEnumTypeFieldNumber:
		enum := descriptor.GetEnumType()[path[1]]
		if !appendEnumToFullName(enum, path, &name) {
			return "", false
		}
	case fileDescriptorProtoServiceFieldNumber:
		service := descriptor.GetService()[path[1]]
		appendNameComponent(service.GetName(), &name)
		if len(path) > 2 {
			if !appendMethodToFullName(service, path, &name) {
				return "", false
			}
		}
	default:
		return "", false
	}

	return name.String(), true
}

func appendMethodToFullName(service *descriptorpb.ServiceDescriptorProto, path []int32, name *strings.Builder) bool {
	if len(path) == 4 && path[2] == serviceDescriptorProtoMethodFieldNumber {
		appendFieldComponent(service.GetMethod()[path[3]].GetName(), name)
		return true
	}
	return false
}

func appendToFullName(descriptor *descriptorpb.DescriptorProto, path []int32, name *strings.Builder) bool {
	switch path[0] {
	case descriptorProtoFieldFieldNumber:
		field := descriptor.GetField()[path[1]]
		appendFieldComponent(field.GetName(), name)
		return true
	case descriptorProtoNestedTypeFieldNumber:
		msg := descriptor.GetNestedType()[path[1]]
		return appendMessageToFullName(msg, path, name)
	case descriptorProtoEnumTypeFieldNumber:
		enum := descriptor.GetEnumType()[path[1]]
		return appendEnumToFullName(enum, path, name)
	default:
		return false
	}
}

func appendMessageToFullName(msg *descriptorpb.DescriptorProto, path []int32, name *strings.Builder) bool {
	appendNameComponent(msg.GetName(), name)
	if len(path) > 2 {
		appendToFullName(msg, path[2:], name)
	}
	return true
}

func appendEnumToFullName(enum *descriptorpb.EnumDescriptorProto, path []int32, name *strings.Builder) bool {
	appendNameComponent(enum.GetName(), name)
	if len(path) <= 2 {
		return true
	}
	if path[2] == enumDescriptorProtoValueFieldNumber {
		appendFieldComponent(enum.GetValue()[path[3]].GetName(), name)
		return true
	}
	return false
}

func appendNameComponent(component string, name *strings.Builder) {
	name.WriteByte('.')
	name.WriteString(component)
}

func appendFieldComponent(component string, name *strings.Builder) {
	name.WriteByte('/')
	name.WriteString(component)
}
