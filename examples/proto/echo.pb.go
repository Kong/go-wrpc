// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: examples/proto/echo.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type EchoRPCRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	S string `protobuf:"bytes,1,opt,name=s,proto3" json:"s,omitempty"`
}

func (x *EchoRPCRequest) Reset() {
	*x = EchoRPCRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_proto_echo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRPCRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRPCRequest) ProtoMessage() {}

func (x *EchoRPCRequest) ProtoReflect() protoreflect.Message {
	mi := &file_examples_proto_echo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRPCRequest.ProtoReflect.Descriptor instead.
func (*EchoRPCRequest) Descriptor() ([]byte, []int) {
	return file_examples_proto_echo_proto_rawDescGZIP(), []int{0}
}

func (x *EchoRPCRequest) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

type EchoRPCResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	S string `protobuf:"bytes,1,opt,name=s,proto3" json:"s,omitempty"`
}

func (x *EchoRPCResponse) Reset() {
	*x = EchoRPCResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_proto_echo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRPCResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRPCResponse) ProtoMessage() {}

func (x *EchoRPCResponse) ProtoReflect() protoreflect.Message {
	mi := &file_examples_proto_echo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRPCResponse.ProtoReflect.Descriptor instead.
func (*EchoRPCResponse) Descriptor() ([]byte, []int) {
	return file_examples_proto_echo_proto_rawDescGZIP(), []int{1}
}

func (x *EchoRPCResponse) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

var File_examples_proto_echo_proto protoreflect.FileDescriptor

var file_examples_proto_echo_proto_rawDesc = []byte{
	0x0a, 0x19, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x65, 0x63, 0x68, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x77, 0x72, 0x70,
	0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x1e, 0x0a, 0x0e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x50, 0x43, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73, 0x22,
	0x1f, 0x0a, 0x0f, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x50, 0x43, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73,
	0x32, 0x57, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x4f, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f,
	0x12, 0x22, 0x2e, 0x77, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x50, 0x43, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x77, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x50,
	0x43, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6f, 0x6e, 0x67, 0x2f, 0x67, 0x6f, 0x2d,
	0x77, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examples_proto_echo_proto_rawDescOnce sync.Once
	file_examples_proto_echo_proto_rawDescData = file_examples_proto_echo_proto_rawDesc
)

func file_examples_proto_echo_proto_rawDescGZIP() []byte {
	file_examples_proto_echo_proto_rawDescOnce.Do(func() {
		file_examples_proto_echo_proto_rawDescData = protoimpl.X.CompressGZIP(file_examples_proto_echo_proto_rawDescData)
	})
	return file_examples_proto_echo_proto_rawDescData
}

var file_examples_proto_echo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_examples_proto_echo_proto_goTypes = []interface{}{
	(*EchoRPCRequest)(nil),  // 0: wrpc.example.proto.EchoRPCRequest
	(*EchoRPCResponse)(nil), // 1: wrpc.example.proto.EchoRPCResponse
}
var file_examples_proto_echo_proto_depIdxs = []int32{
	0, // 0: wrpc.example.proto.Echo.Echo:input_type -> wrpc.example.proto.EchoRPCRequest
	1, // 1: wrpc.example.proto.Echo.Echo:output_type -> wrpc.example.proto.EchoRPCResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_examples_proto_echo_proto_init() }
func file_examples_proto_echo_proto_init() {
	if File_examples_proto_echo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examples_proto_echo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoRPCRequest); i {
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
		file_examples_proto_echo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoRPCResponse); i {
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
			RawDescriptor: file_examples_proto_echo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examples_proto_echo_proto_goTypes,
		DependencyIndexes: file_examples_proto_echo_proto_depIdxs,
		MessageInfos:      file_examples_proto_echo_proto_msgTypes,
	}.Build()
	File_examples_proto_echo_proto = out.File
	file_examples_proto_echo_proto_rawDesc = nil
	file_examples_proto_echo_proto_goTypes = nil
	file_examples_proto_echo_proto_depIdxs = nil
}
