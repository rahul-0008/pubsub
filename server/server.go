package server

import (
	"fmt"
	"log"
	"time"

	pb "github.com/pubsub/proto"
	"google.golang.org/grpc"
)

func RegisterServer(config string, mode string, masterAddr string, port int, opts ...grpc.ServerOption) *grpc.Server {

	server := grpc.NewServer(opts...)

	switch config {
	case "single":
		pb.RegisterPubsubServer(server, NewSingleMachine())

	case "masterSlave":
		switch mode {
		case "master":
			master := NewMasterServer()

			pb.RegisterPubsubServer(server, master)
			pb.RegisterMasterServiceServer(server, master)
			log.Println("Initialising Master Server at ", port)

		case "slave":
			masterAdderss := fmt.Sprintf("%s:%d", masterAddr, port)
			connTimeout := time.Duration(10 * time.Second)

			log.Println("Connectiong to the master")
			conn, err := grpc.Dial(masterAdderss,
				grpc.WithTimeout(connTimeout),
				grpc.WithInsecure())

			if err != nil {
				log.Fatalln("Could not connect to Master, ", err)
				return nil
			}
			slave := NewSlaveServer(conn)
			pb.RegisterMasterServiceServer(server, slave)
			pb.RegisterPubsubServer(server, slave)
		}
	}

	return server
}
