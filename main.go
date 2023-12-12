package main

import (
	"fmt"
	"log"
	"net"
	"time"

	server "github.com/pubsub/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	fmt.Println("Initialising the pubsub")

	srv := server.RegisterServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}))

	l, e := net.Listen("tcp", "0.0.0.0:7476")

	if e != nil {
		log.Println("Failed to listen")
	}

	srv.Serve(l)
}
