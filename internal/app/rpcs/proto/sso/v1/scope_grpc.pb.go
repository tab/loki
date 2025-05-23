// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: sso/v1/scope.proto

package ssov1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ScopeService_List_FullMethodName   = "/sso.v1.ScopeService/List"
	ScopeService_Get_FullMethodName    = "/sso.v1.ScopeService/Get"
	ScopeService_Create_FullMethodName = "/sso.v1.ScopeService/Create"
	ScopeService_Update_FullMethodName = "/sso.v1.ScopeService/Update"
	ScopeService_Delete_FullMethodName = "/sso.v1.ScopeService/Delete"
)

// ScopeServiceClient is the client API for ScopeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Scope service provides CRUD operations for managing scopes
type ScopeServiceClient interface {
	List(ctx context.Context, in *PaginatedListRequest, opts ...grpc.CallOption) (*ListScopesResponse, error)
	Get(ctx context.Context, in *GetScopeRequest, opts ...grpc.CallOption) (*GetScopeResponse, error)
	Create(ctx context.Context, in *CreateScopeRequest, opts ...grpc.CallOption) (*CreateScopeResponse, error)
	Update(ctx context.Context, in *UpdateScopeRequest, opts ...grpc.CallOption) (*UpdateScopeResponse, error)
	Delete(ctx context.Context, in *DeleteScopeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type scopeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScopeServiceClient(cc grpc.ClientConnInterface) ScopeServiceClient {
	return &scopeServiceClient{cc}
}

func (c *scopeServiceClient) List(ctx context.Context, in *PaginatedListRequest, opts ...grpc.CallOption) (*ListScopesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListScopesResponse)
	err := c.cc.Invoke(ctx, ScopeService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scopeServiceClient) Get(ctx context.Context, in *GetScopeRequest, opts ...grpc.CallOption) (*GetScopeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetScopeResponse)
	err := c.cc.Invoke(ctx, ScopeService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scopeServiceClient) Create(ctx context.Context, in *CreateScopeRequest, opts ...grpc.CallOption) (*CreateScopeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateScopeResponse)
	err := c.cc.Invoke(ctx, ScopeService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scopeServiceClient) Update(ctx context.Context, in *UpdateScopeRequest, opts ...grpc.CallOption) (*UpdateScopeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateScopeResponse)
	err := c.cc.Invoke(ctx, ScopeService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scopeServiceClient) Delete(ctx context.Context, in *DeleteScopeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ScopeService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScopeServiceServer is the server API for ScopeService service.
// All implementations must embed UnimplementedScopeServiceServer
// for forward compatibility.
//
// Scope service provides CRUD operations for managing scopes
type ScopeServiceServer interface {
	List(context.Context, *PaginatedListRequest) (*ListScopesResponse, error)
	Get(context.Context, *GetScopeRequest) (*GetScopeResponse, error)
	Create(context.Context, *CreateScopeRequest) (*CreateScopeResponse, error)
	Update(context.Context, *UpdateScopeRequest) (*UpdateScopeResponse, error)
	Delete(context.Context, *DeleteScopeRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedScopeServiceServer()
}

// UnimplementedScopeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedScopeServiceServer struct{}

func (UnimplementedScopeServiceServer) List(context.Context, *PaginatedListRequest) (*ListScopesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedScopeServiceServer) Get(context.Context, *GetScopeRequest) (*GetScopeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedScopeServiceServer) Create(context.Context, *CreateScopeRequest) (*CreateScopeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedScopeServiceServer) Update(context.Context, *UpdateScopeRequest) (*UpdateScopeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedScopeServiceServer) Delete(context.Context, *DeleteScopeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedScopeServiceServer) mustEmbedUnimplementedScopeServiceServer() {}
func (UnimplementedScopeServiceServer) testEmbeddedByValue()                      {}

// UnsafeScopeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScopeServiceServer will
// result in compilation errors.
type UnsafeScopeServiceServer interface {
	mustEmbedUnimplementedScopeServiceServer()
}

func RegisterScopeServiceServer(s grpc.ServiceRegistrar, srv ScopeServiceServer) {
	// If the following call pancis, it indicates UnimplementedScopeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ScopeService_ServiceDesc, srv)
}

func _ScopeService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaginatedListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScopeServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScopeService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScopeServiceServer).List(ctx, req.(*PaginatedListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScopeService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScopeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScopeServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScopeService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScopeServiceServer).Get(ctx, req.(*GetScopeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScopeService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateScopeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScopeServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScopeService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScopeServiceServer).Create(ctx, req.(*CreateScopeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScopeService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateScopeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScopeServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScopeService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScopeServiceServer).Update(ctx, req.(*UpdateScopeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScopeService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteScopeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScopeServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScopeService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScopeServiceServer).Delete(ctx, req.(*DeleteScopeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ScopeService_ServiceDesc is the grpc.ServiceDesc for ScopeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScopeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sso.v1.ScopeService",
	HandlerType: (*ScopeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _ScopeService_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ScopeService_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ScopeService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ScopeService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ScopeService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso/v1/scope.proto",
}
