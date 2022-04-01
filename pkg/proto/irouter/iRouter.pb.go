// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: iRouter.proto

package irouter

import (
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

type ASRoute_Operation int32

const (
	ASRoute_UPDATE     ASRoute_Operation = 0
	ASRoute_WITHDRAWAL ASRoute_Operation = 1
)

// Enum value maps for ASRoute_Operation.
var (
	ASRoute_Operation_name = map[int32]string{
		0: "UPDATE",
		1: "WITHDRAWAL",
	}
	ASRoute_Operation_value = map[string]int32{
		"UPDATE":     0,
		"WITHDRAWAL": 1,
	}
)

func (x ASRoute_Operation) Enum() *ASRoute_Operation {
	p := new(ASRoute_Operation)
	*p = x
	return p
}

func (x ASRoute_Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ASRoute_Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_iRouter_proto_enumTypes[0].Descriptor()
}

func (ASRoute_Operation) Type() protoreflect.EnumType {
	return &file_iRouter_proto_enumTypes[0]
}

func (x ASRoute_Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ASRoute_Operation.Descriptor instead.
func (ASRoute_Operation) EnumDescriptor() ([]byte, []int) {
	return file_iRouter_proto_rawDescGZIP(), []int{0, 0}
}

type ASRoute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peer          string            `protobuf:"bytes,1,opt,name=Peer,proto3" json:"Peer,omitempty"`
	PeerAS        string            `protobuf:"bytes,2,opt,name=PeerAS,proto3" json:"PeerAS,omitempty"`
	Host          string            `protobuf:"bytes,3,opt,name=Host,proto3" json:"Host,omitempty"`
	Op            ASRoute_Operation `protobuf:"varint,4,opt,name=Op,proto3,enum=irouter.ASRoute_Operation" json:"Op,omitempty"`
	Path          string            `protobuf:"bytes,5,opt,name=Path,proto3" json:"Path,omitempty"`
	Community     string            `protobuf:"bytes,6,opt,name=Community,proto3" json:"Community,omitempty"`
	Origin        string            `protobuf:"bytes,7,opt,name=Origin,proto3" json:"Origin,omitempty"`
	Announcements string            `protobuf:"bytes,8,opt,name=Announcements,proto3" json:"Announcements,omitempty"`
	Raw           []byte            `protobuf:"bytes,9,opt,name=Raw,proto3" json:"Raw,omitempty"`
}

func (x *ASRoute) Reset() {
	*x = ASRoute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iRouter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ASRoute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ASRoute) ProtoMessage() {}

func (x *ASRoute) ProtoReflect() protoreflect.Message {
	mi := &file_iRouter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ASRoute.ProtoReflect.Descriptor instead.
func (*ASRoute) Descriptor() ([]byte, []int) {
	return file_iRouter_proto_rawDescGZIP(), []int{0}
}

func (x *ASRoute) GetPeer() string {
	if x != nil {
		return x.Peer
	}
	return ""
}

func (x *ASRoute) GetPeerAS() string {
	if x != nil {
		return x.PeerAS
	}
	return ""
}

func (x *ASRoute) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *ASRoute) GetOp() ASRoute_Operation {
	if x != nil {
		return x.Op
	}
	return ASRoute_UPDATE
}

func (x *ASRoute) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ASRoute) GetCommunity() string {
	if x != nil {
		return x.Community
	}
	return ""
}

func (x *ASRoute) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *ASRoute) GetAnnouncements() string {
	if x != nil {
		return x.Announcements
	}
	return ""
}

func (x *ASRoute) GetRaw() []byte {
	if x != nil {
		return x.Raw
	}
	return nil
}

var File_iRouter_proto protoreflect.FileDescriptor

var file_iRouter_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x69, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x69, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x22, 0xa0, 0x02, 0x0a, 0x07, 0x41, 0x53, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x65, 0x65, 0x72,
	0x41, 0x53, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x50, 0x65, 0x65, 0x72, 0x41, 0x53,
	0x12, 0x12, 0x0a, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x48, 0x6f, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x02, 0x4f, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1a, 0x2e, 0x69, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x41, 0x53, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x02, 0x4f, 0x70,
	0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74,
	0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x41, 0x6e,
	0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x52, 0x61, 0x77, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x52,
	0x61, 0x77, 0x22, 0x27, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0a, 0x0a, 0x06, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x57,
	0x49, 0x54, 0x48, 0x44, 0x52, 0x41, 0x57, 0x41, 0x4c, 0x10, 0x01, 0x42, 0x0f, 0x5a, 0x0d, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iRouter_proto_rawDescOnce sync.Once
	file_iRouter_proto_rawDescData = file_iRouter_proto_rawDesc
)

func file_iRouter_proto_rawDescGZIP() []byte {
	file_iRouter_proto_rawDescOnce.Do(func() {
		file_iRouter_proto_rawDescData = protoimpl.X.CompressGZIP(file_iRouter_proto_rawDescData)
	})
	return file_iRouter_proto_rawDescData
}

var file_iRouter_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_iRouter_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_iRouter_proto_goTypes = []interface{}{
	(ASRoute_Operation)(0), // 0: irouter.ASRoute.Operation
	(*ASRoute)(nil),        // 1: irouter.ASRoute
}
var file_iRouter_proto_depIdxs = []int32{
	0, // 0: irouter.ASRoute.Op:type_name -> irouter.ASRoute.Operation
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_iRouter_proto_init() }
func file_iRouter_proto_init() {
	if File_iRouter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_iRouter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ASRoute); i {
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
			RawDescriptor: file_iRouter_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iRouter_proto_goTypes,
		DependencyIndexes: file_iRouter_proto_depIdxs,
		EnumInfos:         file_iRouter_proto_enumTypes,
		MessageInfos:      file_iRouter_proto_msgTypes,
	}.Build()
	File_iRouter_proto = out.File
	file_iRouter_proto_rawDesc = nil
	file_iRouter_proto_goTypes = nil
	file_iRouter_proto_depIdxs = nil
}
