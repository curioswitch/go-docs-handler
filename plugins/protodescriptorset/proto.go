package protodescriptorset

// Go does not generate constants for field numbers, so we vendor
// in what we need as they won't change.

const (
	fileDescriptorProtoMessageTypeFieldNumber = 4
	fileDescriptorProtoEnumTypeFieldNumber    = 5
	fileDescriptorProtoServiceFieldNumber     = 6
)

const serviceDescriptorProtoMethodFieldNumber = 2

const (
	descriptorProtoFieldFieldNumber      = 2
	descriptorProtoNestedTypeFieldNumber = 3
	descriptorProtoEnumTypeFieldNumber   = 4
)

const enumDescriptorProtoValueFieldNumber = 2
