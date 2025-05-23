// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: sso/v1/role.proto

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

// Role represents a role object
type Role struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	PermissionIds []string               `protobuf:"bytes,4,rep,name=permission_ids,json=permissionIds,proto3" json:"permission_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Role) Reset() {
	*x = Role{}
	mi := &file_sso_v1_role_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{0}
}

func (x *Role) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Role) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Role) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Role) GetPermissionIds() []string {
	if x != nil {
		return x.PermissionIds
	}
	return nil
}

// ListRolesResponse is the response for the List method
type ListRolesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*Role                `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	Meta          *PaginationMeta        `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRolesResponse) Reset() {
	*x = ListRolesResponse{}
	mi := &file_sso_v1_role_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRolesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRolesResponse) ProtoMessage() {}

func (x *ListRolesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRolesResponse.ProtoReflect.Descriptor instead.
func (*ListRolesResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{1}
}

func (x *ListRolesResponse) GetData() []*Role {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ListRolesResponse) GetMeta() *PaginationMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

// GetRoleRequest is the request for the Get method
type GetRoleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRoleRequest) Reset() {
	*x = GetRoleRequest{}
	mi := &file_sso_v1_role_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleRequest) ProtoMessage() {}

func (x *GetRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleRequest.ProtoReflect.Descriptor instead.
func (*GetRoleRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{2}
}

func (x *GetRoleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetRoleResponse is the response for the Get method
type GetRoleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Role                  `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRoleResponse) Reset() {
	*x = GetRoleResponse{}
	mi := &file_sso_v1_role_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleResponse) ProtoMessage() {}

func (x *GetRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleResponse.ProtoReflect.Descriptor instead.
func (*GetRoleResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{3}
}

func (x *GetRoleResponse) GetData() *Role {
	if x != nil {
		return x.Data
	}
	return nil
}

// CreateRoleRequest is the request for the Create method
type CreateRoleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	PermissionIds []string               `protobuf:"bytes,3,rep,name=permission_ids,json=permissionIds,proto3" json:"permission_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateRoleRequest) Reset() {
	*x = CreateRoleRequest{}
	mi := &file_sso_v1_role_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRoleRequest) ProtoMessage() {}

func (x *CreateRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRoleRequest.ProtoReflect.Descriptor instead.
func (*CreateRoleRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{4}
}

func (x *CreateRoleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRoleRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateRoleRequest) GetPermissionIds() []string {
	if x != nil {
		return x.PermissionIds
	}
	return nil
}

// CreateRoleResponse is the response for the Create method
type CreateRoleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Role                  `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateRoleResponse) Reset() {
	*x = CreateRoleResponse{}
	mi := &file_sso_v1_role_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRoleResponse) ProtoMessage() {}

func (x *CreateRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRoleResponse.ProtoReflect.Descriptor instead.
func (*CreateRoleResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{5}
}

func (x *CreateRoleResponse) GetData() *Role {
	if x != nil {
		return x.Data
	}
	return nil
}

// UpdateRoleRequest is the request for the Update method
type UpdateRoleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	PermissionIds []string               `protobuf:"bytes,4,rep,name=permission_ids,json=permissionIds,proto3" json:"permission_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRoleRequest) Reset() {
	*x = UpdateRoleRequest{}
	mi := &file_sso_v1_role_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRoleRequest) ProtoMessage() {}

func (x *UpdateRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRoleRequest.ProtoReflect.Descriptor instead.
func (*UpdateRoleRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateRoleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateRoleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateRoleRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateRoleRequest) GetPermissionIds() []string {
	if x != nil {
		return x.PermissionIds
	}
	return nil
}

// UpdateRoleResponse is the response for the Update method
type UpdateRoleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Role                  `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRoleResponse) Reset() {
	*x = UpdateRoleResponse{}
	mi := &file_sso_v1_role_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRoleResponse) ProtoMessage() {}

func (x *UpdateRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRoleResponse.ProtoReflect.Descriptor instead.
func (*UpdateRoleResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateRoleResponse) GetData() *Role {
	if x != nil {
		return x.Data
	}
	return nil
}

// DeleteRoleRequest is the request for the Delete method
type DeleteRoleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteRoleRequest) Reset() {
	*x = DeleteRoleRequest{}
	mi := &file_sso_v1_role_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRoleRequest) ProtoMessage() {}

func (x *DeleteRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_role_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRoleRequest.ProtoReflect.Descriptor instead.
func (*DeleteRoleRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_role_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteRoleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_sso_v1_role_proto protoreflect.FileDescriptor

const file_sso_v1_role_proto_rawDesc = "" +
	"\n" +
	"\x11sso/v1/role.proto\x12\x06sso.v1\x1a\x1bbuf/validate/validate.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x17sso/v1/pagination.proto\"\xa3\x01\n" +
	"\x04Role\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12\x1d\n" +
	"\x04name\x18\x02 \x01(\tB\t\xbaH\x06r\x04\x10\x01\x18dR\x04name\x12,\n" +
	"\vdescription\x18\x03 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\x80\x10R\vdescription\x124\n" +
	"\x0epermission_ids\x18\x04 \x03(\tB\r\xbaH\n" +
	"\x92\x01\a\"\x05r\x03\xb0\x01\x01R\rpermissionIds\"a\n" +
	"\x11ListRolesResponse\x12 \n" +
	"\x04data\x18\x01 \x03(\v2\f.sso.v1.RoleR\x04data\x12*\n" +
	"\x04meta\x18\x02 \x01(\v2\x16.sso.v1.PaginationMetaR\x04meta\"*\n" +
	"\x0eGetRoleRequest\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id\"3\n" +
	"\x0fGetRoleResponse\x12 \n" +
	"\x04data\x18\x01 \x01(\v2\f.sso.v1.RoleR\x04data\"\x96\x01\n" +
	"\x11CreateRoleRequest\x12\x1d\n" +
	"\x04name\x18\x01 \x01(\tB\t\xbaH\x06r\x04\x10\x01\x18dR\x04name\x12,\n" +
	"\vdescription\x18\x02 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\x80\x10R\vdescription\x124\n" +
	"\x0epermission_ids\x18\x03 \x03(\tB\r\xbaH\n" +
	"\x92\x01\a\"\x05r\x03\xb0\x01\x01R\rpermissionIds\"6\n" +
	"\x12CreateRoleResponse\x12 \n" +
	"\x04data\x18\x01 \x01(\v2\f.sso.v1.RoleR\x04data\"\xb0\x01\n" +
	"\x11UpdateRoleRequest\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12\x1d\n" +
	"\x04name\x18\x02 \x01(\tB\t\xbaH\x06r\x04\x10\x01\x18dR\x04name\x12,\n" +
	"\vdescription\x18\x03 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\x80\x10R\vdescription\x124\n" +
	"\x0epermission_ids\x18\x04 \x03(\tB\r\xbaH\n" +
	"\x92\x01\a\"\x05r\x03\xb0\x01\x01R\rpermissionIds\"6\n" +
	"\x12UpdateRoleResponse\x12 \n" +
	"\x04data\x18\x01 \x01(\v2\f.sso.v1.RoleR\x04data\"-\n" +
	"\x11DeleteRoleRequest\x12\x18\n" +
	"\x02id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x02id2\xcf\x02\n" +
	"\vRoleService\x12A\n" +
	"\x04List\x12\x1c.sso.v1.PaginatedListRequest\x1a\x19.sso.v1.ListRolesResponse\"\x00\x128\n" +
	"\x03Get\x12\x16.sso.v1.GetRoleRequest\x1a\x17.sso.v1.GetRoleResponse\"\x00\x12A\n" +
	"\x06Create\x12\x19.sso.v1.CreateRoleRequest\x1a\x1a.sso.v1.CreateRoleResponse\"\x00\x12A\n" +
	"\x06Update\x12\x19.sso.v1.UpdateRoleRequest\x1a\x1a.sso.v1.UpdateRoleResponse\"\x00\x12=\n" +
	"\x06Delete\x12\x19.sso.v1.DeleteRoleRequest\x1a\x16.google.protobuf.Empty\"\x00B+Z)loki/internal/app/rpcs/proto/sso/v1;ssov1b\x06proto3"

var (
	file_sso_v1_role_proto_rawDescOnce sync.Once
	file_sso_v1_role_proto_rawDescData []byte
)

func file_sso_v1_role_proto_rawDescGZIP() []byte {
	file_sso_v1_role_proto_rawDescOnce.Do(func() {
		file_sso_v1_role_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sso_v1_role_proto_rawDesc), len(file_sso_v1_role_proto_rawDesc)))
	})
	return file_sso_v1_role_proto_rawDescData
}

var file_sso_v1_role_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_sso_v1_role_proto_goTypes = []any{
	(*Role)(nil),                 // 0: sso.v1.Role
	(*ListRolesResponse)(nil),    // 1: sso.v1.ListRolesResponse
	(*GetRoleRequest)(nil),       // 2: sso.v1.GetRoleRequest
	(*GetRoleResponse)(nil),      // 3: sso.v1.GetRoleResponse
	(*CreateRoleRequest)(nil),    // 4: sso.v1.CreateRoleRequest
	(*CreateRoleResponse)(nil),   // 5: sso.v1.CreateRoleResponse
	(*UpdateRoleRequest)(nil),    // 6: sso.v1.UpdateRoleRequest
	(*UpdateRoleResponse)(nil),   // 7: sso.v1.UpdateRoleResponse
	(*DeleteRoleRequest)(nil),    // 8: sso.v1.DeleteRoleRequest
	(*PaginationMeta)(nil),       // 9: sso.v1.PaginationMeta
	(*PaginatedListRequest)(nil), // 10: sso.v1.PaginatedListRequest
	(*emptypb.Empty)(nil),        // 11: google.protobuf.Empty
}
var file_sso_v1_role_proto_depIdxs = []int32{
	0,  // 0: sso.v1.ListRolesResponse.data:type_name -> sso.v1.Role
	9,  // 1: sso.v1.ListRolesResponse.meta:type_name -> sso.v1.PaginationMeta
	0,  // 2: sso.v1.GetRoleResponse.data:type_name -> sso.v1.Role
	0,  // 3: sso.v1.CreateRoleResponse.data:type_name -> sso.v1.Role
	0,  // 4: sso.v1.UpdateRoleResponse.data:type_name -> sso.v1.Role
	10, // 5: sso.v1.RoleService.List:input_type -> sso.v1.PaginatedListRequest
	2,  // 6: sso.v1.RoleService.Get:input_type -> sso.v1.GetRoleRequest
	4,  // 7: sso.v1.RoleService.Create:input_type -> sso.v1.CreateRoleRequest
	6,  // 8: sso.v1.RoleService.Update:input_type -> sso.v1.UpdateRoleRequest
	8,  // 9: sso.v1.RoleService.Delete:input_type -> sso.v1.DeleteRoleRequest
	1,  // 10: sso.v1.RoleService.List:output_type -> sso.v1.ListRolesResponse
	3,  // 11: sso.v1.RoleService.Get:output_type -> sso.v1.GetRoleResponse
	5,  // 12: sso.v1.RoleService.Create:output_type -> sso.v1.CreateRoleResponse
	7,  // 13: sso.v1.RoleService.Update:output_type -> sso.v1.UpdateRoleResponse
	11, // 14: sso.v1.RoleService.Delete:output_type -> google.protobuf.Empty
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_sso_v1_role_proto_init() }
func file_sso_v1_role_proto_init() {
	if File_sso_v1_role_proto != nil {
		return
	}
	file_sso_v1_pagination_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sso_v1_role_proto_rawDesc), len(file_sso_v1_role_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sso_v1_role_proto_goTypes,
		DependencyIndexes: file_sso_v1_role_proto_depIdxs,
		MessageInfos:      file_sso_v1_role_proto_msgTypes,
	}.Build()
	File_sso_v1_role_proto = out.File
	file_sso_v1_role_proto_goTypes = nil
	file_sso_v1_role_proto_depIdxs = nil
}
