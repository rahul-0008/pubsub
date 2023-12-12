// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: proto/pubsub.proto

package proto

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

// PubsubClient is the client API for Pubsub service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PubsubClient interface {
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Pubsub_SubscribeClient, error)
}

type pubsubClient struct {
	cc grpc.ClientConnInterface
}

func NewPubsubClient(cc grpc.ClientConnInterface) PubsubClient {
	return &pubsubClient{cc}
}

func (c *pubsubClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, "/proto.pubsub/publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubsubClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Pubsub_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Pubsub_ServiceDesc.Streams[0], "/proto.pubsub/subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &pubsubSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Pubsub_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type pubsubSubscribeClient struct {
	grpc.ClientStream
}

func (x *pubsubSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PubsubServer is the server API for Pubsub service.
// All implementations must embed UnimplementedPubsubServer
// for forward compatibility
type PubsubServer interface {
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	Subscribe(*SubscribeRequest, Pubsub_SubscribeServer) error
	mustEmbedUnimplementedPubsubServer()
}

// UnimplementedPubsubServer must be embedded to have forward compatible implementations.
type UnimplementedPubsubServer struct {
}

func (UnimplementedPubsubServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedPubsubServer) Subscribe(*SubscribeRequest, Pubsub_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedPubsubServer) mustEmbedUnimplementedPubsubServer() {}

// UnsafePubsubServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PubsubServer will
// result in compilation errors.
type UnsafePubsubServer interface {
	mustEmbedUnimplementedPubsubServer()
}

func RegisterPubsubServer(s grpc.ServiceRegistrar, srv PubsubServer) {
	s.RegisterService(&Pubsub_ServiceDesc, srv)
}

func _Pubsub_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubsubServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.pubsub/publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubsubServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pubsub_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PubsubServer).Subscribe(m, &pubsubSubscribeServer{stream})
}

type Pubsub_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type pubsubSubscribeServer struct {
	grpc.ServerStream
}

func (x *pubsubSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Pubsub_ServiceDesc is the grpc.ServiceDesc for Pubsub service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pubsub_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.pubsub",
	HandlerType: (*PubsubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "publish",
			Handler:    _Pubsub_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "subscribe",
			Handler:       _Pubsub_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/pubsub.proto",
}