package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/subcommands"
	pb "github.com/pubsub/proto"
	"google.golang.org/grpc"
)

var (
	address = flag.String("address", "localhost", "pubsub server address")
	port    = flag.Int("port", 7476, "pubsub server port")
	timeout = flag.String("timeout", "10s", "pubsub server connection time out")
)

func createServeiceStub(server string) pb.PubsubClient {
	connTimeout, err := time.ParseDuration(*timeout)
	if err != nil {
		log.Fatalf("Invalid timeout format : %s", *timeout)
	}
	conn, err := grpc.Dial(server, grpc.WithTimeout(connTimeout), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("cannot connect to the server [%s] : %v", server, err)
	}

	return pb.NewPubsubClient(conn)
}

func main() {
	flag.Parse()
	server := fmt.Sprintf("%s:%d", *address, *port)
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&PublishCmd{server: server}, "")
	subcommands.Register(&subscribeCmd{server: server}, "")

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
