// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: idl/demo/demo.proto

package demo

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

// DemoServiceClient is the client API for DemoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DemoServiceClient interface {
	// Echo 样例接口
	Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error)
}

type demoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDemoServiceClient(cc grpc.ClientConnInterface) DemoServiceClient {
	return &demoServiceClient{cc}
}

func (c *demoServiceClient) Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := c.cc.Invoke(ctx, "/demo.DemoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DemoServiceServer is the server API for DemoService service.
// All implementations must embed UnimplementedDemoServiceServer
// for forward compatibility
type DemoServiceServer interface {
	// Echo 样例接口
	Echo(context.Context, *StringMessage) (*StringMessage, error)
	mustEmbedUnimplementedDemoServiceServer()
}

// UnimplementedDemoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDemoServiceServer struct {
}

func (UnimplementedDemoServiceServer) Echo(context.Context, *StringMessage) (*StringMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedDemoServiceServer) mustEmbedUnimplementedDemoServiceServer() {}

// UnsafeDemoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DemoServiceServer will
// result in compilation errors.
type UnsafeDemoServiceServer interface {
	mustEmbedUnimplementedDemoServiceServer()
}

func RegisterDemoServiceServer(s grpc.ServiceRegistrar, srv DemoServiceServer) {
	s.RegisterService(&DemoService_ServiceDesc, srv)
}

func _DemoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DemoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.DemoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DemoServiceServer).Echo(ctx, req.(*StringMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// DemoService_ServiceDesc is the grpc.ServiceDesc for DemoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DemoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "demo.DemoService",
	HandlerType: (*DemoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _DemoService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/demo/demo.proto",
}
