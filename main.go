package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	server "github.com/pubsub/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	config        = flag.String("config", "single", "Cinfigurationmode")
	mode          = flag.String("mode", "master", "Mode")
	masterAdderss = flag.String("address", "0.0.0.0", "Master Address")
	masterPort    = flag.Int("masterPort", 7476, "Port")
	port          = flag.Int("port", 7476, "Port")
)

func main() {
	flag.Parse()
	fmt.Println("Initialising the pubsub in mode: ", *config)

	srv := server.RegisterServer(*config, *mode, *masterAdderss, *masterPort, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}))

	address := fmt.Sprintf("%s:%d", *masterAdderss, *port)
	l, e := net.Listen("tcp", address)

	if e != nil {
		log.Println("Failed to listen")
	}

	srv.Serve(l)
}
