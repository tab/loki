// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: sso/v1/scope.proto

package ssov1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Scope represents a scope object
type Scope struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Scope) Reset() {
	*x = Scope{}
	mi := &file_sso_v1_scope_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Scope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Scope) ProtoMessage() {}

func (x *Scope) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Scope.ProtoReflect.Descriptor instead.
func (*Scope) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{0}
}

func (x *Scope) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Scope) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Scope) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// ListScopesResponse is the response for the List method
type ListScopesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*Scope               `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	Meta          *PaginationMeta        `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListScopesResponse) Reset() {
	*x = ListScopesResponse{}
	mi := &file_sso_v1_scope_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListScopesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListScopesResponse) ProtoMessage() {}

func (x *ListScopesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListScopesResponse.ProtoReflect.Descriptor instead.
func (*ListScopesResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{1}
}

func (x *ListScopesResponse) GetData() []*Scope {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ListScopesResponse) GetMeta() *PaginationMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

// GetScopeRequest is the request for the Get method
type GetScopeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScopeRequest) Reset() {
	*x = GetScopeRequest{}
	mi := &file_sso_v1_scope_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScopeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScopeRequest) ProtoMessage() {}

func (x *GetScopeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScopeRequest.ProtoReflect.Descriptor instead.
func (*GetScopeRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{2}
}

func (x *GetScopeRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetScopeResponse is the response for the Get method
type GetScopeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Scope                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScopeResponse) Reset() {
	*x = GetScopeResponse{}
	mi := &file_sso_v1_scope_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScopeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScopeResponse) ProtoMessage() {}

func (x *GetScopeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScopeResponse.ProtoReflect.Descriptor instead.
func (*GetScopeResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{3}
}

func (x *GetScopeResponse) GetData() *Scope {
	if x != nil {
		return x.Data
	}
	return nil
}

// CreateScopeRequest is the request for the Create method
type CreateScopeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateScopeRequest) Reset() {
	*x = CreateScopeRequest{}
	mi := &file_sso_v1_scope_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateScopeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScopeRequest) ProtoMessage() {}

func (x *CreateScopeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScopeRequest.ProtoReflect.Descriptor instead.
func (*CreateScopeRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{4}
}

func (x *CreateScopeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateScopeRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// CreateScopeResponse is the response for the Create method
type CreateScopeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Scope                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateScopeResponse) Reset() {
	*x = CreateScopeResponse{}
	mi := &file_sso_v1_scope_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateScopeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScopeResponse) ProtoMessage() {}

func (x *CreateScopeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScopeResponse.ProtoReflect.Descriptor instead.
func (*CreateScopeResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{5}
}

func (x *CreateScopeResponse) GetData() *Scope {
	if x != nil {
		return x.Data
	}
	return nil
}

// UpdateScopeRequest is the request for the Update method
type UpdateScopeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateScopeRequest) Reset() {
	*x = UpdateScopeRequest{}
	mi := &file_sso_v1_scope_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateScopeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateScopeRequest) ProtoMessage() {}

func (x *UpdateScopeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateScopeRequest.ProtoReflect.Descriptor instead.
func (*UpdateScopeRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateScopeRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateScopeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateScopeRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// UpdateScopeResponse is the response for the Update method
type UpdateScopeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Scope                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateScopeResponse) Reset() {
	*x = UpdateScopeResponse{}
	mi := &file_sso_v1_scope_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateScopeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateScopeResponse) ProtoMessage() {}

func (x *UpdateScopeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateScopeResponse.ProtoReflect.Descriptor instead.
func (*UpdateScopeResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateScopeResponse) GetData() *Scope {
	if x != nil {
		return x.Data
	}
	return nil
}

// DeleteScopeRequest is the request for the Delete method
type DeleteScopeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteScopeRequest) Reset() {
	*x = DeleteScopeRequest{}
	mi := &file_sso_v1_scope_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteScopeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteScopeRequest) ProtoMessage() {}

func (x *DeleteScopeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_scope_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteScopeRequest.ProtoReflect.Descriptor instead.
func (*DeleteScopeRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_scope_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteScopeRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_sso_v1_scope_proto protoreflect.FileDescriptor

const file_sso_v1_scope_proto_rawDesc = "" +
	"\n" +
	"\x12sso/v1/scope.proto\x12\x06sso.v1\x1a\x1bbuf/validate/validate.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x17sso/v1/pagination.proto\"n\n" +
	"\x05Scope\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12\x1d\n" +
	"\x04name\x18\x02 \x01(\tB\t\xbaH\x06r\x04\x10\x01\x18dR\x04name\x12,\n" +
	"\vdescription\x18\x03 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\x80\x10R\vdescription\"c\n" +
	"\x12ListScopesResponse\x12!\n" +
	"\x04data\x18\x01 \x03(\v2\r.sso.v1.ScopeR\x04data\x12*\n" +
	"\x04meta\x18\x02 \x01(\v2\x16.sso.v1.PaginationMetaR\x04meta\"+\n" +
	"\x0fGetScopeRequest\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id\"5\n" +
	"\x10GetScopeResponse\x12!\n" +
	"\x04data\x18\x01 \x01(\v2\r.sso.v1.ScopeR\x04data\"a\n" +
	"\x12CreateScopeRequest\x12\x1d\n" +
	"\x04name\x18\x01 \x01(\tB\t\xbaH\x06r\x04\x10\x01\x18dR\x04name\x12,\n" +
	"\vdescription\x18\x02 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\x80\x10R\vdescription\"8\n" +
	"\x13CreateScopeResponse\x12!\n" +
	"\x04data\x18\x01 \x01(\v2\r.sso.v1.ScopeR\x04data\"{\n" +
	"\x12UpdateScopeRequest\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12\x1d\n" +
	"\x04name\x18\x02 \x01(\tB\t\xbaH\x06r\x04\x10\x01\x18dR\x04name\x12,\n" +
	"\vdescription\x18\x03 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\x80\x10R\vdescription\"8\n" +
	"\x13UpdateScopeResponse\x12!\n" +
	"\x04data\x18\x01 \x01(\v2\r.sso.v1.ScopeR\x04data\".\n" +
	"\x12DeleteScopeRequest\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id2\xd8\x02\n" +
	"\fScopeService\x12B\n" +
	"\x04List\x12\x1c.sso.v1.PaginatedListRequest\x1a\x1a.sso.v1.ListScopesResponse\"\x00\x12:\n" +
	"\x03Get\x12\x17.sso.v1.GetScopeRequest\x1a\x18.sso.v1.GetScopeResponse\"\x00\x12C\n" +
	"\x06Create\x12\x1a.sso.v1.CreateScopeRequest\x1a\x1b.sso.v1.CreateScopeResponse\"\x00\x12C\n" +
	"\x06Update\x12\x1a.sso.v1.UpdateScopeRequest\x1a\x1b.sso.v1.UpdateScopeResponse\"\x00\x12>\n" +
	"\x06Delete\x12\x1a.sso.v1.DeleteScopeRequest\x1a\x16.google.protobuf.Empty\"\x00B+Z)loki/internal/app/rpcs/proto/sso/v1;ssov1b\x06proto3"

var (
	file_sso_v1_scope_proto_rawDescOnce sync.Once
	file_sso_v1_scope_proto_rawDescData []byte
)

func file_sso_v1_scope_proto_rawDescGZIP() []byte {
	file_sso_v1_scope_proto_rawDescOnce.Do(func() {
		file_sso_v1_scope_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sso_v1_scope_proto_rawDesc), len(file_sso_v1_scope_proto_rawDesc)))
	})
	return file_sso_v1_scope_proto_rawDescData
}

var file_sso_v1_scope_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_sso_v1_scope_proto_goTypes = []any{
	(*Scope)(nil),                // 0: sso.v1.Scope
	(*ListScopesResponse)(nil),   // 1: sso.v1.ListScopesResponse
	(*GetScopeRequest)(nil),      // 2: sso.v1.GetScopeRequest
	(*GetScopeResponse)(nil),     // 3: sso.v1.GetScopeResponse
	(*CreateScopeRequest)(nil),   // 4: sso.v1.CreateScopeRequest
	(*CreateScopeResponse)(nil),  // 5: sso.v1.CreateScopeResponse
	(*UpdateScopeRequest)(nil),   // 6: sso.v1.UpdateScopeRequest
	(*UpdateScopeResponse)(nil),  // 7: sso.v1.UpdateScopeResponse
	(*DeleteScopeRequest)(nil),   // 8: sso.v1.DeleteScopeRequest
	(*PaginationMeta)(nil),       // 9: sso.v1.PaginationMeta
	(*PaginatedListRequest)(nil), // 10: sso.v1.PaginatedListRequest
	(*emptypb.Empty)(nil),        // 11: google.protobuf.Empty
}
var file_sso_v1_scope_proto_depIdxs = []int32{
	0,  // 0: sso.v1.ListScopesResponse.data:type_name -> sso.v1.Scope
	9,  // 1: sso.v1.ListScopesResponse.meta:type_name -> sso.v1.PaginationMeta
	0,  // 2: sso.v1.GetScopeResponse.data:type_name -> sso.v1.Scope
	0,  // 3: sso.v1.CreateScopeResponse.data:type_name -> sso.v1.Scope
	0,  // 4: sso.v1.UpdateScopeResponse.data:type_name -> sso.v1.Scope
	10, // 5: sso.v1.ScopeService.List:input_type -> sso.v1.PaginatedListRequest
	2,  // 6: sso.v1.ScopeService.Get:input_type -> sso.v1.GetScopeRequest
	4,  // 7: sso.v1.ScopeService.Create:input_type -> sso.v1.CreateScopeRequest
	6,  // 8: sso.v1.ScopeService.Update:input_type -> sso.v1.UpdateScopeRequest
	8,  // 9: sso.v1.ScopeService.Delete:input_type -> sso.v1.DeleteScopeRequest
	1,  // 10: sso.v1.ScopeService.List:output_type -> sso.v1.ListScopesResponse
	3,  // 11: sso.v1.ScopeService.Get:output_type -> sso.v1.GetScopeResponse
	5,  // 12: sso.v1.ScopeService.Create:output_type -> sso.v1.CreateScopeResponse
	7,  // 13: sso.v1.ScopeService.Update:output_type -> sso.v1.UpdateScopeResponse
	11, // 14: sso.v1.ScopeService.Delete:output_type -> google.protobuf.Empty
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_sso_v1_scope_proto_init() }
func file_sso_v1_scope_proto_init() {
	if File_sso_v1_scope_proto != nil {
		return
	}
	file_sso_v1_pagination_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sso_v1_scope_proto_rawDesc), len(file_sso_v1_scope_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sso_v1_scope_proto_goTypes,
		DependencyIndexes: file_sso_v1_scope_proto_depIdxs,
		MessageInfos:      file_sso_v1_scope_proto_msgTypes,
	}.Build()
	File_sso_v1_scope_proto = out.File
	file_sso_v1_scope_proto_goTypes = nil
	file_sso_v1_scope_proto_depIdxs = nil
}
