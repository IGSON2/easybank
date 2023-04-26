// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: service_easy_bank.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EasybankClient is the client API for Easybank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EasybankClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type easybankClient struct {
	cc grpc.ClientConnInterface
}

func NewEasybankClient(cc grpc.ClientConnInterface) EasybankClient {
	return &easybankClient{cc}
}

func (c *easybankClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/pb.Easybank/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EasybankServer is the server API for Easybank service.
// All implementations must embed UnimplementedEasybankServer
// for forward compatibility
type EasybankServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	mustEmbedUnimplementedEasybankServer()
}

// UnimplementedEasybankServer must be embedded to have forward compatible implementations.
type UnimplementedEasybankServer struct {
}

func (UnimplementedEasybankServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedEasybankServer) mustEmbedUnimplementedEasybankServer() {}

// UnsafeEasybankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EasybankServer will
// result in compilation errors.
type UnsafeEasybankServer interface {
	mustEmbedUnimplementedEasybankServer()
}

func RegisterEasybankServer(s grpc.ServiceRegistrar, srv EasybankServer) {
	s.RegisterService(&Easybank_ServiceDesc, srv)
}

func _Easybank_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EasybankServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Easybank/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EasybankServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Easybank_ServiceDesc is the grpc.ServiceDesc for Easybank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Easybank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Easybank",
	HandlerType: (*EasybankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _Easybank_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_easy_bank.proto",
}
