// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.27.0
// source: user.proto

package userPb

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

// LSRequest 登录或注册请求,其中User可以表示用户名或邮箱
type LSRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User     string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Phone    string `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	Captcha  string `protobuf:"bytes,4,opt,name=captcha,proto3" json:"captcha,omitempty"`
	Ip       string `protobuf:"bytes,5,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *LSRequest) Reset() {
	*x = LSRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LSRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LSRequest) ProtoMessage() {}

func (x *LSRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LSRequest.ProtoReflect.Descriptor instead.
func (*LSRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *LSRequest) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *LSRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *LSRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *LSRequest) GetCaptcha() string {
	if x != nil {
		return x.Captcha
	}
	return ""
}

func (x *LSRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    *LoginResponse_Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	UserInfo *User                `protobuf:"bytes,3,opt,name=UserInfo,proto3" json:"UserInfo,omitempty"`
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetToken() *LoginResponse_Token {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *LoginResponse) GetUserInfo() *User {
	if x != nil {
		return x.UserInfo
	}
	return nil
}

type EmptyLSResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyLSResponse) Reset() {
	*x = EmptyLSResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyLSResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyLSResponse) ProtoMessage() {}

func (x *EmptyLSResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyLSResponse.ProtoReflect.Descriptor instead.
func (*EmptyLSResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

type GetUserInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId int64 `protobuf:"varint,1,opt,name=actorId,proto3" json:"actorId,omitempty"`
	UserId  int64 `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetUserInfoRequest) Reset() {
	*x = GetUserInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoRequest) ProtoMessage() {}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoRequest.ProtoReflect.Descriptor instead.
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserInfoRequest) GetActorId() int64 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *GetUserInfoRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetUserInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserInfoResponse) Reset() {
	*x = GetUserInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoResponse) ProtoMessage() {}

func (x *GetUserInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoResponse.ProtoReflect.Descriptor instead.
func (*GetUserInfoResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserInfoResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId           int64   `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	TotalLiked       *int64  `protobuf:"varint,2,opt,name=totalLiked,proto3,oneof" json:"totalLiked,omitempty"`
	LikeCount        *int64  `protobuf:"varint,3,opt,name=likeCount,proto3,oneof" json:"likeCount,omitempty"`
	CollectCount     *int64  `protobuf:"varint,4,opt,name=collectCount,proto3,oneof" json:"collectCount,omitempty"`
	PostCount        *int64  `protobuf:"varint,5,opt,name=postCount,proto3,oneof" json:"postCount,omitempty"`
	FollowCount      *int64  `protobuf:"varint,6,opt,name=followCount,proto3,oneof" json:"followCount,omitempty"`
	FansCount        *int64  `protobuf:"varint,7,opt,name=fansCount,proto3,oneof" json:"fansCount,omitempty"`
	Username         string  `protobuf:"bytes,8,opt,name=username,proto3" json:"username,omitempty"`
	Email            *string `protobuf:"bytes,9,opt,name=email,proto3,oneof" json:"email,omitempty"`
	Avatar           *string `protobuf:"bytes,10,opt,name=avatar,proto3,oneof" json:"avatar,omitempty"`
	Birthday         *string `protobuf:"bytes,11,opt,name=birthday,proto3,oneof" json:"birthday,omitempty"`
	Introduction     *string `protobuf:"bytes,12,opt,name=introduction,proto3,oneof" json:"introduction,omitempty"`
	Phone            *string `protobuf:"bytes,13,opt,name=phone,proto3,oneof" json:"phone,omitempty"`
	School           *string `protobuf:"bytes,14,opt,name=school,proto3,oneof" json:"school,omitempty"`
	LastLoginTime    *string `protobuf:"bytes,15,opt,name=lastLoginTime,proto3,oneof" json:"lastLoginTime,omitempty"`
	LastLoginIp      *string `protobuf:"bytes,16,opt,name=lastLoginIp,proto3,oneof" json:"lastLoginIp,omitempty"`
	NoticeInfo       *string `protobuf:"bytes,17,opt,name=noticeInfo,proto3,oneof" json:"noticeInfo,omitempty"`
	JoinTime         *string `protobuf:"bytes,18,opt,name=joinTime,proto3,oneof" json:"joinTime,omitempty"`
	TotalCoinCount   *uint32 `protobuf:"varint,19,opt,name=totalCoinCount,proto3,oneof" json:"totalCoinCount,omitempty"`
	CurrentCoinCount *uint32 `protobuf:"varint,20,opt,name=currentCoinCount,proto3,oneof" json:"currentCoinCount,omitempty"`
	Theme            *uint32 `protobuf:"varint,21,opt,name=theme,proto3,oneof" json:"theme,omitempty"`
	Sex              *uint32 `protobuf:"varint,22,opt,name=sex,proto3,oneof" json:"sex,omitempty"`
	Status           *uint32 `protobuf:"varint,23,opt,name=status,proto3,oneof" json:"status,omitempty"`
	IsFollow         bool    `protobuf:"varint,24,opt,name=isFollow,proto3" json:"isFollow,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *User) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *User) GetTotalLiked() int64 {
	if x != nil && x.TotalLiked != nil {
		return *x.TotalLiked
	}
	return 0
}

func (x *User) GetLikeCount() int64 {
	if x != nil && x.LikeCount != nil {
		return *x.LikeCount
	}
	return 0
}

func (x *User) GetCollectCount() int64 {
	if x != nil && x.CollectCount != nil {
		return *x.CollectCount
	}
	return 0
}

func (x *User) GetPostCount() int64 {
	if x != nil && x.PostCount != nil {
		return *x.PostCount
	}
	return 0
}

func (x *User) GetFollowCount() int64 {
	if x != nil && x.FollowCount != nil {
		return *x.FollowCount
	}
	return 0
}

func (x *User) GetFansCount() int64 {
	if x != nil && x.FansCount != nil {
		return *x.FansCount
	}
	return 0
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *User) GetAvatar() string {
	if x != nil && x.Avatar != nil {
		return *x.Avatar
	}
	return ""
}

func (x *User) GetBirthday() string {
	if x != nil && x.Birthday != nil {
		return *x.Birthday
	}
	return ""
}

func (x *User) GetIntroduction() string {
	if x != nil && x.Introduction != nil {
		return *x.Introduction
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil && x.Phone != nil {
		return *x.Phone
	}
	return ""
}

func (x *User) GetSchool() string {
	if x != nil && x.School != nil {
		return *x.School
	}
	return ""
}

func (x *User) GetLastLoginTime() string {
	if x != nil && x.LastLoginTime != nil {
		return *x.LastLoginTime
	}
	return ""
}

func (x *User) GetLastLoginIp() string {
	if x != nil && x.LastLoginIp != nil {
		return *x.LastLoginIp
	}
	return ""
}

func (x *User) GetNoticeInfo() string {
	if x != nil && x.NoticeInfo != nil {
		return *x.NoticeInfo
	}
	return ""
}

func (x *User) GetJoinTime() string {
	if x != nil && x.JoinTime != nil {
		return *x.JoinTime
	}
	return ""
}

func (x *User) GetTotalCoinCount() uint32 {
	if x != nil && x.TotalCoinCount != nil {
		return *x.TotalCoinCount
	}
	return 0
}

func (x *User) GetCurrentCoinCount() uint32 {
	if x != nil && x.CurrentCoinCount != nil {
		return *x.CurrentCoinCount
	}
	return 0
}

func (x *User) GetTheme() uint32 {
	if x != nil && x.Theme != nil {
		return *x.Theme
	}
	return 0
}

func (x *User) GetSex() uint32 {
	if x != nil && x.Sex != nil {
		return *x.Sex
	}
	return 0
}

func (x *User) GetStatus() uint32 {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return 0
}

func (x *User) GetIsFollow() bool {
	if x != nil {
		return x.IsFollow
	}
	return false
}

type GetUserExistInformationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetUserExistInformationRequest) Reset() {
	*x = GetUserExistInformationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserExistInformationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserExistInformationRequest) ProtoMessage() {}

func (x *GetUserExistInformationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserExistInformationRequest.ProtoReflect.Descriptor instead.
func (*GetUserExistInformationRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{6}
}

func (x *GetUserExistInformationRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetUserExistInformationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Existed bool `protobuf:"varint,1,opt,name=existed,proto3" json:"existed,omitempty"`
}

func (x *GetUserExistInformationResponse) Reset() {
	*x = GetUserExistInformationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserExistInformationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserExistInformationResponse) ProtoMessage() {}

func (x *GetUserExistInformationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserExistInformationResponse.ProtoReflect.Descriptor instead.
func (*GetUserExistInformationResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{7}
}

func (x *GetUserExistInformationResponse) GetExisted() bool {
	if x != nil {
		return x.Existed
	}
	return false
}

type LoginResponse_Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken  string `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refreshToken,proto3" json:"refreshToken,omitempty"`
}

func (x *LoginResponse_Token) Reset() {
	*x = LoginResponse_Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse_Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse_Token) ProtoMessage() {}

func (x *LoginResponse_Token) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse_Token.ProtoReflect.Descriptor instead.
func (*LoginResponse_Token) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1, 0}
}

func (x *LoginResponse_Token) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *LoginResponse_Token) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x50, 0x62, 0x22, 0x7b, 0x0a, 0x09, 0x4c, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x61, 0x70, 0x74, 0x63,
	0x68, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x61, 0x70, 0x74, 0x63, 0x68,
	0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x70, 0x22, 0xbb, 0x01, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x28, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50,
	0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x1a, 0x4d, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x72,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x11, 0x0a, 0x0f, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4c, 0x53, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x46, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x37, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x20, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x22, 0xd8, 0x08, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4c, 0x69, 0x6b,
	0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x4c, 0x69, 0x6b, 0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6c, 0x69, 0x6b,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x09,
	0x6c, 0x69, 0x6b, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0c,
	0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x48, 0x02, 0x52, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x03, 0x52, 0x09, 0x70, 0x6f, 0x73, 0x74,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x66, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x48, 0x04, 0x52,
	0x0b, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12,
	0x21, 0x0a, 0x09, 0x66, 0x61, 0x6e, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x48, 0x05, 0x52, 0x09, 0x66, 0x61, 0x6e, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x07, 0x52, 0x06, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64,
	0x61, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52, 0x08, 0x62, 0x69, 0x72, 0x74,
	0x68, 0x64, 0x61, 0x79, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x09, 0x52,
	0x0c, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x12, 0x19, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x0a, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x73,
	0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0b, 0x52, 0x06, 0x73,
	0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x0c, 0x52, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x49, 0x70, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0d, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x49, 0x70, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x6e, 0x6f,
	0x74, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0e,
	0x52, 0x0a, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x88, 0x01, 0x01, 0x12,
	0x1f, 0x0a, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x0f, 0x52, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x2b, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x10, 0x52, 0x0e, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x43, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x2f, 0x0a,
	0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x11, 0x52, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x43, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x19,
	0x0a, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x12, 0x52,
	0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x73, 0x65, 0x78,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x13, 0x52, 0x03, 0x73, 0x65, 0x78, 0x88, 0x01, 0x01,
	0x12, 0x1b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0d,
	0x48, 0x14, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x1a, 0x0a,
	0x08, 0x69, 0x73, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x69, 0x73, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x4c, 0x69, 0x6b, 0x65, 0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6c, 0x69, 0x6b,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x70, 0x6f, 0x73, 0x74,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x66, 0x61, 0x6e, 0x73, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x62, 0x69, 0x72,
	0x74, 0x68, 0x64, 0x61, 0x79, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x42, 0x10, 0x0a, 0x0e, 0x5f,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x0e, 0x0a,
	0x0c, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x49, 0x70, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x0b, 0x0a, 0x09,
	0x5f, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x43, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x13, 0x0a, 0x11,
	0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f,
	0x73, 0x65, 0x78, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x38,
	0x0a, 0x1e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78, 0x69, 0x73, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3b, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x45, 0x78, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65,
	0x78, 0x69, 0x73, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x78,
	0x69, 0x73, 0x74, 0x65, 0x64, 0x32, 0xf2, 0x02, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x11, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e,
	0x4c, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x50, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x43, 0x61, 0x70, 0x74, 0x63,
	0x68, 0x61, 0x12, 0x11, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x4c, 0x53, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36,
	0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x11, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50,
	0x62, 0x2e, 0x4c, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x50, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4c, 0x53, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6a,
	0x0a, 0x17, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78, 0x69, 0x73, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x50, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78, 0x69, 0x73, 0x74, 0x49,
	0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x27, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x45, 0x78, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1f, 0x5a, 0x1d, 0x73, 0x74,
	0x61, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x50, 0x62, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_user_proto_goTypes = []interface{}{
	(*LSRequest)(nil),                       // 0: userPb.LSRequest
	(*LoginResponse)(nil),                   // 1: userPb.LoginResponse
	(*EmptyLSResponse)(nil),                 // 2: userPb.EmptyLSResponse
	(*GetUserInfoRequest)(nil),              // 3: userPb.GetUserInfoRequest
	(*GetUserInfoResponse)(nil),             // 4: userPb.GetUserInfoResponse
	(*User)(nil),                            // 5: userPb.User
	(*GetUserExistInformationRequest)(nil),  // 6: userPb.GetUserExistInformationRequest
	(*GetUserExistInformationResponse)(nil), // 7: userPb.GetUserExistInformationResponse
	(*LoginResponse_Token)(nil),             // 8: userPb.LoginResponse.Token
}
var file_user_proto_depIdxs = []int32{
	8, // 0: userPb.LoginResponse.token:type_name -> userPb.LoginResponse.Token
	5, // 1: userPb.LoginResponse.UserInfo:type_name -> userPb.User
	5, // 2: userPb.GetUserInfoResponse.user:type_name -> userPb.User
	0, // 3: userPb.userService.LoginPassword:input_type -> userPb.LSRequest
	0, // 4: userPb.userService.LoginCaptcha:input_type -> userPb.LSRequest
	0, // 5: userPb.userService.Signup:input_type -> userPb.LSRequest
	3, // 6: userPb.userService.GetUserInfo:input_type -> userPb.GetUserInfoRequest
	6, // 7: userPb.userService.GetUserExistInformation:input_type -> userPb.GetUserExistInformationRequest
	1, // 8: userPb.userService.LoginPassword:output_type -> userPb.LoginResponse
	1, // 9: userPb.userService.LoginCaptcha:output_type -> userPb.LoginResponse
	2, // 10: userPb.userService.Signup:output_type -> userPb.EmptyLSResponse
	4, // 11: userPb.userService.GetUserInfo:output_type -> userPb.GetUserInfoResponse
	7, // 12: userPb.userService.GetUserExistInformation:output_type -> userPb.GetUserExistInformationResponse
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LSRequest); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginResponse); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyLSResponse); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserInfoRequest); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserInfoResponse); i {
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
		file_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserExistInformationRequest); i {
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
		file_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserExistInformationResponse); i {
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
		file_user_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginResponse_Token); i {
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
	file_user_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
