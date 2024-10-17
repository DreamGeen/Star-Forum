// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.27.0
// source: publish.proto

package publishPb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	feedPb "star/proto/feed/feedPb"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreatePostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      int64  `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	CommunityId int64  `protobuf:"varint,2,opt,name=CommunityId,proto3" json:"CommunityId,omitempty"`
	Content     string `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content,omitempty"`
	IsScan      bool   `protobuf:"varint,5,opt,name=IsScan,proto3" json:"IsScan,omitempty"`
}

func (x *CreatePostRequest) Reset() {
	*x = CreatePostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publish_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostRequest) ProtoMessage() {}

func (x *CreatePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_publish_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostRequest.ProtoReflect.Descriptor instead.
func (*CreatePostRequest) Descriptor() ([]byte, []int) {
	return file_publish_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePostRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreatePostRequest) GetCommunityId() int64 {
	if x != nil {
		return x.CommunityId
	}
	return 0
}

func (x *CreatePostRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *CreatePostRequest) GetIsScan() bool {
	if x != nil {
		return x.IsScan
	}
	return false
}

type CreatePostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreatePostResponse) Reset() {
	*x = CreatePostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publish_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostResponse) ProtoMessage() {}

func (x *CreatePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_publish_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostResponse.ProtoReflect.Descriptor instead.
func (*CreatePostResponse) Descriptor() ([]byte, []int) {
	return file_publish_proto_rawDescGZIP(), []int{1}
}

type CountPostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
}

func (x *CountPostRequest) Reset() {
	*x = CountPostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publish_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CountPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CountPostRequest) ProtoMessage() {}

func (x *CountPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_publish_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CountPostRequest.ProtoReflect.Descriptor instead.
func (*CountPostRequest) Descriptor() ([]byte, []int) {
	return file_publish_proto_rawDescGZIP(), []int{2}
}

func (x *CountPostRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type CountPostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int64 `protobuf:"varint,1,opt,name=Count,proto3" json:"Count,omitempty"`
}

func (x *CountPostResponse) Reset() {
	*x = CountPostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publish_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CountPostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CountPostResponse) ProtoMessage() {}

func (x *CountPostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_publish_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CountPostResponse.ProtoReflect.Descriptor instead.
func (*CountPostResponse) Descriptor() ([]byte, []int) {
	return file_publish_proto_rawDescGZIP(), []int{3}
}

func (x *CountPostResponse) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ListPostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId int64 `protobuf:"varint,1,opt,name=ActorId,proto3" json:"ActorId,omitempty"`
	UserId  int64 `protobuf:"varint,2,opt,name=UserId,proto3" json:"UserId,omitempty"`
}

func (x *ListPostRequest) Reset() {
	*x = ListPostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publish_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPostRequest) ProtoMessage() {}

func (x *ListPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_publish_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPostRequest.ProtoReflect.Descriptor instead.
func (*ListPostRequest) Descriptor() ([]byte, []int) {
	return file_publish_proto_rawDescGZIP(), []int{4}
}

func (x *ListPostRequest) GetActorId() int64 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *ListPostRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ListPostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Posts []*feedPb.Post `protobuf:"bytes,1,rep,name=Posts,proto3" json:"Posts,omitempty"`
}

func (x *ListPostResponse) Reset() {
	*x = ListPostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publish_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPostResponse) ProtoMessage() {}

func (x *ListPostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_publish_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPostResponse.ProtoReflect.Descriptor instead.
func (*ListPostResponse) Descriptor() ([]byte, []int) {
	return file_publish_proto_rawDescGZIP(), []int{5}
}

func (x *ListPostResponse) GetPosts() []*feedPb.Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

var File_publish_proto protoreflect.FileDescriptor

var file_publish_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x1a, 0x1a, 0x73, 0x74, 0x61, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x65, 0x64, 0x2f, 0x66, 0x65, 0x65, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x49, 0x73, 0x53, 0x63, 0x61, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x49, 0x73, 0x53, 0x63, 0x61, 0x6e, 0x22, 0x14, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2a, 0x0a,
	0x10, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x29, 0x0a, 0x11, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x43, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x63, 0x74, 0x6f, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x10, 0x4c, 0x69, 0x73,
	0x74, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a,
	0x05, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x66,
	0x65, 0x65, 0x64, 0x50, 0x62, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x05, 0x50, 0x6f, 0x73, 0x74,
	0x73, 0x32, 0xe8, 0x01, 0x0a, 0x0e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f,
	0x73, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x46, 0x0a, 0x09, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x50, 0x62, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x6f, 0x73, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x28, 0x5a, 0x26,
	0x73, 0x74, 0x61, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x3b, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x50, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_publish_proto_rawDescOnce sync.Once
	file_publish_proto_rawDescData = file_publish_proto_rawDesc
)

func file_publish_proto_rawDescGZIP() []byte {
	file_publish_proto_rawDescOnce.Do(func() {
		file_publish_proto_rawDescData = protoimpl.X.CompressGZIP(file_publish_proto_rawDescData)
	})
	return file_publish_proto_rawDescData
}

var file_publish_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_publish_proto_goTypes = []interface{}{
	(*CreatePostRequest)(nil),  // 0: publishPb.CreatePostRequest
	(*CreatePostResponse)(nil), // 1: publishPb.CreatePostResponse
	(*CountPostRequest)(nil),   // 2: publishPb.CountPostRequest
	(*CountPostResponse)(nil),  // 3: publishPb.CountPostResponse
	(*ListPostRequest)(nil),    // 4: publishPb.ListPostRequest
	(*ListPostResponse)(nil),   // 5: publishPb.ListPostResponse
	(*feedPb.Post)(nil),        // 6: feedPb.Post
}
var file_publish_proto_depIdxs = []int32{
	6, // 0: publishPb.ListPostResponse.Posts:type_name -> feedPb.Post
	0, // 1: publishPb.PublishService.CreatePost:input_type -> publishPb.CreatePostRequest
	2, // 2: publishPb.PublishService.CountPost:input_type -> publishPb.CountPostRequest
	4, // 3: publishPb.PublishService.ListPost:input_type -> publishPb.ListPostRequest
	1, // 4: publishPb.PublishService.CreatePost:output_type -> publishPb.CreatePostResponse
	3, // 5: publishPb.PublishService.CountPost:output_type -> publishPb.CountPostResponse
	5, // 6: publishPb.PublishService.ListPost:output_type -> publishPb.ListPostResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_publish_proto_init() }
func file_publish_proto_init() {
	if File_publish_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_publish_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePostRequest); i {
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
		file_publish_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePostResponse); i {
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
		file_publish_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CountPostRequest); i {
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
		file_publish_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CountPostResponse); i {
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
		file_publish_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPostRequest); i {
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
		file_publish_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPostResponse); i {
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
			RawDescriptor: file_publish_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_publish_proto_goTypes,
		DependencyIndexes: file_publish_proto_depIdxs,
		MessageInfos:      file_publish_proto_msgTypes,
	}.Build()
	File_publish_proto = out.File
	file_publish_proto_rawDesc = nil
	file_publish_proto_goTypes = nil
	file_publish_proto_depIdxs = nil
}
