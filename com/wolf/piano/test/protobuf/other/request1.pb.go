// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.12.1
// source: gopiano/com/wolf/piano/test/protobuf/other/request1.proto

// package 属于 proto 文件自身的范围定义，与生成的 go 代码无关
// 为了避免当导入其他 proto 文件时导致的文件内的命名冲突。所以，当导入非本包的 message 时，需要加 package 前缀,同包内的引用不需要加包名前缀。

package other

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
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

var File_gopiano_com_wolf_piano_test_protobuf_other_request1_proto protoreflect.FileDescriptor

var file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_rawDesc = []byte{
	0x0a, 0x39, 0x67, 0x6f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6f,
	0x6c, 0x66, 0x2f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x2f, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6f, 0x74, 0x68,
	0x65, 0x72, 0x1a, 0x3b, 0x67, 0x6f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x6f, 0x6c, 0x66, 0x2f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x2f, 0x6e,
	0x65, 0x77, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42,
	0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69,
	0x79, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x63, 0x6f, 0x6d,
	0x2f, 0x77, 0x6f, 0x6c, 0x66, 0x2f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x50,
	0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_goTypes = []interface{}{}
var file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_init() }
func file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_init() {
	if File_gopiano_com_wolf_piano_test_protobuf_other_request1_proto != nil {
		return
	}
	file_gopiano_com_wolf_piano_test_protobuf_other_newrequest_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_goTypes,
		DependencyIndexes: file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_depIdxs,
	}.Build()
	File_gopiano_com_wolf_piano_test_protobuf_other_request1_proto = out.File
	file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_rawDesc = nil
	file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_goTypes = nil
	file_gopiano_com_wolf_piano_test_protobuf_other_request1_proto_depIdxs = nil
}