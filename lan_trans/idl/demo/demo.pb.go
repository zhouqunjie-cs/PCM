// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: idl/demo/demo.proto

package demo

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OurTeam int32

const (
	// github: devad
	OurTeam_devad OurTeam = 0
)

// Enum value maps for OurTeam.
var (
	OurTeam_name = map[int32]string{
		0: "devad",
	}
	OurTeam_value = map[string]int32{
		"devad": 0,
	}
)

func (x OurTeam) Enum() *OurTeam {
	p := new(OurTeam)
	*p = x
	return p
}

func (x OurTeam) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OurTeam) Descriptor() protoreflect.EnumDescriptor {
	return file_idl_demo_demo_proto_enumTypes[0].Descriptor()
}

func (OurTeam) Type() protoreflect.EnumType {
	return &file_idl_demo_demo_proto_enumTypes[0]
}

func (x OurTeam) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OurTeam.Descriptor instead.
func (OurTeam) EnumDescriptor() ([]byte, []int) {
	return file_idl_demo_demo_proto_rawDescGZIP(), []int{0}
}

type StringMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringMessage) Reset() {
	*x = StringMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_idl_demo_demo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringMessage) ProtoMessage() {}

func (x *StringMessage) ProtoReflect() protoreflect.Message {
	mi := &file_idl_demo_demo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringMessage.ProtoReflect.Descriptor instead.
func (*StringMessage) Descriptor() ([]byte, []int) {
	return file_idl_demo_demo_proto_rawDescGZIP(), []int{0}
}

func (x *StringMessage) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_idl_demo_demo_proto protoreflect.FileDescriptor

var file_idl_demo_demo_proto_rawDesc = []byte{
	0x0a, 0x13, 0x69, 0x64, 0x6c, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x65, 0x6d, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x0d, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x2a, 0x14, 0x0a, 0x07, 0x4f, 0x75, 0x72, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x09, 0x0a, 0x05, 0x64,
	0x65, 0x76, 0x61, 0x64, 0x10, 0x00, 0x32, 0xb3, 0x01, 0x0a, 0x0b, 0x44, 0x65, 0x6d, 0x6f, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xa3, 0x01, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12,
	0x13, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x13, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x71, 0x92, 0x41, 0x59, 0x22, 0x53,
	0x0a, 0x21, 0x46, 0x69, 0x6e, 0x64, 0x20, 0x6f, 0x75, 0x74, 0x20, 0x6d, 0x6f, 0x72, 0x65, 0x20,
	0x61, 0x62, 0x6f, 0x75, 0x74, 0x20, 0x74, 0x68, 0x65, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x12, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x65, 0x63, 0x6f,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x58, 0x01, 0x62, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x22, 0x0a, 0x2f,
	0x61, 0x70, 0x69, 0x73, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x3a, 0x01, 0x2a, 0x42, 0x30, 0x5a, 0x2e,
	0x67, 0x69, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x6e, 0x2f, 0x4a,
	0x43, 0x43, 0x45, 0x2f, 0x50, 0x43, 0x4d, 0x2f, 0x6c, 0x61, 0x6e, 0x5f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x64, 0x6c, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_idl_demo_demo_proto_rawDescOnce sync.Once
	file_idl_demo_demo_proto_rawDescData = file_idl_demo_demo_proto_rawDesc
)

func file_idl_demo_demo_proto_rawDescGZIP() []byte {
	file_idl_demo_demo_proto_rawDescOnce.Do(func() {
		file_idl_demo_demo_proto_rawDescData = protoimpl.X.CompressGZIP(file_idl_demo_demo_proto_rawDescData)
	})
	return file_idl_demo_demo_proto_rawDescData
}

var file_idl_demo_demo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_idl_demo_demo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_idl_demo_demo_proto_goTypes = []interface{}{
	(OurTeam)(0),          // 0: demo.OurTeam
	(*StringMessage)(nil), // 1: demo.StringMessage
}
var file_idl_demo_demo_proto_depIdxs = []int32{
	1, // 0: demo.DemoService.Echo:input_type -> demo.StringMessage
	1, // 1: demo.DemoService.Echo:output_type -> demo.StringMessage
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_idl_demo_demo_proto_init() }
func file_idl_demo_demo_proto_init() {
	if File_idl_demo_demo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_idl_demo_demo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringMessage); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_idl_demo_demo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_idl_demo_demo_proto_goTypes,
		DependencyIndexes: file_idl_demo_demo_proto_depIdxs,
		EnumInfos:         file_idl_demo_demo_proto_enumTypes,
		MessageInfos:      file_idl_demo_demo_proto_msgTypes,
	}.Build()
	File_idl_demo_demo_proto = out.File
	file_idl_demo_demo_proto_rawDesc = nil
	file_idl_demo_demo_proto_goTypes = nil
	file_idl_demo_demo_proto_depIdxs = nil
}
