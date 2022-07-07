// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: payload.proto

package types

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

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	// Publish sends a stream of payloads to the server for in-memory storage
	Publish(ctx context.Context, opts ...grpc.CallOption) (Broker_PublishClient, error)
	// Consume consumes a specified number of objects from the server
	Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Broker_ConsumeClient, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) Publish(ctx context.Context, opts ...grpc.CallOption) (Broker_PublishClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broker_ServiceDesc.Streams[0], "/types.Broker/Publish", opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerPublishClient{stream}
	return x, nil
}

type Broker_PublishClient interface {
	Send(*Payload) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type brokerPublishClient struct {
	grpc.ClientStream
}

func (x *brokerPublishClient) Send(m *Payload) error {
	return x.ClientStream.SendMsg(m)
}

func (x *brokerPublishClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *brokerClient) Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Broker_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broker_ServiceDesc.Streams[1], "/types.Broker/Consume", opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerConsumeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Broker_ConsumeClient interface {
	Recv() (*Payload, error)
	grpc.ClientStream
}

type brokerConsumeClient struct {
	grpc.ClientStream
}

func (x *brokerConsumeClient) Recv() (*Payload, error) {
	m := new(Payload)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility
type BrokerServer interface {
	// Publish sends a stream of payloads to the server for in-memory storage
	Publish(Broker_PublishServer) error
	// Consume consumes a specified number of objects from the server
	Consume(*ConsumeRequest, Broker_ConsumeServer) error
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServer struct {
}

func (UnimplementedBrokerServer) Publish(Broker_PublishServer) error {
	return status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedBrokerServer) Consume(*ConsumeRequest, Broker_ConsumeServer) error {
	return status.Errorf(codes.Unimplemented, "method Consume not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_Publish_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BrokerServer).Publish(&brokerPublishServer{stream})
}

type Broker_PublishServer interface {
	SendAndClose(*Response) error
	Recv() (*Payload, error)
	grpc.ServerStream
}

type brokerPublishServer struct {
	grpc.ServerStream
}

func (x *brokerPublishServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *brokerPublishServer) Recv() (*Payload, error) {
	m := new(Payload)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Broker_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BrokerServer).Consume(m, &brokerConsumeServer{stream})
}

type Broker_ConsumeServer interface {
	Send(*Payload) error
	grpc.ServerStream
}

type brokerConsumeServer struct {
	grpc.ServerStream
}

func (x *brokerConsumeServer) Send(m *Payload) error {
	return x.ServerStream.SendMsg(m)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "types.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Publish",
			Handler:       _Broker_Publish_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Consume",
			Handler:       _Broker_Consume_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "payload.proto",
}
