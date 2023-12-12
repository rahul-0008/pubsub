package server

import (
	pb "github.com/pubsub/proto"
	"google.golang.org/grpc"
)

func RegisterServer(opts ...grpc.ServerOption) *grpc.Server {

	server := grpc.NewServer(opts...)

	pb.RegisterPubsubServer(server, NewSingleMachine())
	return server

}
