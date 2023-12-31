// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.1
// source: greet.proto

package greet

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A name that is known so we can transmit it more efficiently.
type Greeting_KnownName int32

const (
	// An unknown known name. Wait, what does that mean?
	Greeting_UNKNOWN Greeting_KnownName = 0
	// Alice's name.
	Greeting_ALICE Greeting_KnownName = 1
	// Bob's name
	Greeting_BOB Greeting_KnownName = 2
)

// Enum value maps for Greeting_KnownName.
var (
	Greeting_KnownName_name = map[int32]string{
		0: "UNKNOWN",
		1: "ALICE",
		2: "BOB",
	}
	Greeting_KnownName_value = map[string]int32{
		"UNKNOWN": 0,
		"ALICE":   1,
		"BOB":     2,
	}
)

func (x Greeting_KnownName) Enum() *Greeting_KnownName {
	p := new(Greeting_KnownName)
	*p = x
	return p
}

func (x Greeting_KnownName) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Greeting_KnownName) Descriptor() protoreflect.EnumDescriptor {
	return file_greet_proto_enumTypes[0].Descriptor()
}

func (Greeting_KnownName) Type() protoreflect.EnumType {
	return &file_greet_proto_enumTypes[0]
}

func (x Greeting_KnownName) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Greeting_KnownName.Descriptor instead.
func (Greeting_KnownName) EnumDescriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{0, 0}
}

// A greeting to a particular person.
type Greeting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A name that can be used to greet a person.
	//
	// Types that are assignable to Name:
	//
	//	*Greeting_Nickname
	//	*Greeting_FullName_
	//	*Greeting_KnownName_
	Name isGreeting_Name `protobuf_oneof:"name"`
}

func (x *Greeting) Reset() {
	*x = Greeting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_greet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Greeting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Greeting) ProtoMessage() {}

func (x *Greeting) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Greeting.ProtoReflect.Descriptor instead.
func (*Greeting) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{0}
}

func (m *Greeting) GetName() isGreeting_Name {
	if m != nil {
		return m.Name
	}
	return nil
}

func (x *Greeting) GetNickname() string {
	if x, ok := x.GetName().(*Greeting_Nickname); ok {
		return x.Nickname
	}
	return ""
}

func (x *Greeting) GetFullName() *Greeting_FullName {
	if x, ok := x.GetName().(*Greeting_FullName_); ok {
		return x.FullName
	}
	return nil
}

func (x *Greeting) GetKnownName() Greeting_KnownName {
	if x, ok := x.GetName().(*Greeting_KnownName_); ok {
		return x.KnownName
	}
	return Greeting_UNKNOWN
}

type isGreeting_Name interface {
	isGreeting_Name()
}

type Greeting_Nickname struct {
	// A simple nickname.
	Nickname string `protobuf:"bytes,1,opt,name=nickname,proto3,oneof"`
}

type Greeting_FullName_ struct {
	// A full name.
	FullName *Greeting_FullName `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3,oneof"`
}

type Greeting_KnownName_ struct {
	// A well known name.
	KnownName Greeting_KnownName `protobuf:"varint,3,opt,name=known_name,json=knownName,proto3,enum=greet.Greeting_KnownName,oneof"`
}

func (*Greeting_Nickname) isGreeting_Name() {}

func (*Greeting_FullName_) isGreeting_Name() {}

func (*Greeting_KnownName_) isGreeting_Name() {}

// A request to perform a greeting.
type GreetingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The greeting to perform.
	Greeting *Greeting `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
}

func (x *GreetingRequest) Reset() {
	*x = GreetingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_greet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GreetingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GreetingRequest) ProtoMessage() {}

func (x *GreetingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GreetingRequest.ProtoReflect.Descriptor instead.
func (*GreetingRequest) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{1}
}

func (x *GreetingRequest) GetGreeting() *Greeting {
	if x != nil {
		return x.Greeting
	}
	return nil
}

// A response to a greeting.
type GreetingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The result of the greeting.
	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *GreetingResponse) Reset() {
	*x = GreetingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_greet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GreetingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GreetingResponse) ProtoMessage() {}

func (x *GreetingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GreetingResponse.ProtoReflect.Descriptor instead.
func (*GreetingResponse) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{2}
}

func (x *GreetingResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

// A full name of a person.
type Greeting_FullName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The first name.
	FirstName string `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	// The last name.
	LastName string `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
}

func (x *Greeting_FullName) Reset() {
	*x = Greeting_FullName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_greet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Greeting_FullName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Greeting_FullName) ProtoMessage() {}

func (x *Greeting_FullName) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Greeting_FullName.ProtoReflect.Descriptor instead.
func (*Greeting_FullName) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Greeting_FullName) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *Greeting_FullName) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

var File_greet_proto protoreflect.FileDescriptor

var file_greet_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x72, 0x65, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x67,
	0x72, 0x65, 0x65, 0x74, 0x22, 0x9b, 0x02, 0x0a, 0x08, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x12, 0x1c, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x37, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x2e, 0x47, 0x72, 0x65, 0x65, 0x74,
	0x69, 0x6e, 0x67, 0x2e, 0x46, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x48, 0x00, 0x52, 0x08,
	0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x0a, 0x6b, 0x6e, 0x6f, 0x77,
	0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x67,
	0x72, 0x65, 0x65, 0x74, 0x2e, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4b, 0x6e,
	0x6f, 0x77, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x48, 0x00, 0x52, 0x09, 0x6b, 0x6e, 0x6f, 0x77, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x46, 0x0a, 0x08, 0x46, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2c, 0x0a, 0x09,
	0x4b, 0x6e, 0x6f, 0x77, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x4c, 0x49, 0x43, 0x45, 0x10,
	0x01, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x4f, 0x42, 0x10, 0x02, 0x42, 0x06, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0x3e, 0x0a, 0x0f, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x2e,
	0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x22, 0x2a, 0x0a, 0x10, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x4a,
	0x0a, 0x0c, 0x47, 0x72, 0x65, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a,
	0x0a, 0x05, 0x47, 0x72, 0x65, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x2e,
	0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x2e, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x75, 0x72, 0x69, 0x6f, 0x73, 0x77,
	0x69, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2d, 0x64, 0x6f, 0x63, 0x73, 0x2d, 0x68, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x67, 0x72, 0x65, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_greet_proto_rawDescOnce sync.Once
	file_greet_proto_rawDescData = file_greet_proto_rawDesc
)

func file_greet_proto_rawDescGZIP() []byte {
	file_greet_proto_rawDescOnce.Do(func() {
		file_greet_proto_rawDescData = protoimpl.X.CompressGZIP(file_greet_proto_rawDescData)
	})
	return file_greet_proto_rawDescData
}

var file_greet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_greet_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_greet_proto_goTypes = []interface{}{
	(Greeting_KnownName)(0),   // 0: greet.Greeting.KnownName
	(*Greeting)(nil),          // 1: greet.Greeting
	(*GreetingRequest)(nil),   // 2: greet.GreetingRequest
	(*GreetingResponse)(nil),  // 3: greet.GreetingResponse
	(*Greeting_FullName)(nil), // 4: greet.Greeting.FullName
}
var file_greet_proto_depIdxs = []int32{
	4, // 0: greet.Greeting.full_name:type_name -> greet.Greeting.FullName
	0, // 1: greet.Greeting.known_name:type_name -> greet.Greeting.KnownName
	1, // 2: greet.GreetingRequest.greeting:type_name -> greet.Greeting
	2, // 3: greet.GreetService.Greet:input_type -> greet.GreetingRequest
	3, // 4: greet.GreetService.Greet:output_type -> greet.GreetingResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_greet_proto_init() }
func file_greet_proto_init() {
	if File_greet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_greet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Greeting); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_greet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GreetingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_greet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GreetingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_greet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Greeting_FullName); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_greet_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Greeting_Nickname)(nil),
		(*Greeting_FullName_)(nil),
		(*Greeting_KnownName_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_greet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_greet_proto_goTypes,
		DependencyIndexes: file_greet_proto_depIdxs,
		EnumInfos:         file_greet_proto_enumTypes,
		MessageInfos:      file_greet_proto_msgTypes,
	}.Build()
	File_greet_proto = out.File
	file_greet_proto_rawDesc = nil
	file_greet_proto_goTypes = nil
	file_greet_proto_depIdxs = nil
}
