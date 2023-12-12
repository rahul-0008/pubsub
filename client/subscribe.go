package main

import (
	"context"
	"flag"
	"io"
	"log"

	"github.com/google/subcommands"
	pb "github.com/pubsub/proto"
)

type subscribeCmd struct {
	server string
}

func (s *subscribeCmd) Name() string {
	return "subscribe"
}
func (s *subscribeCmd) Synopsis() string {
	return "Subscribe command to subscribe a topic mentioned"
}

func (s *subscribeCmd) Usage() string {
	return `subscribe --topic=<topic> <topic>

	Publish message to the topic
  
  `

}
func (s *subscribeCmd) SetFlags(f *flag.FlagSet) {
}

func (s *subscribeCmd) Execute(ctxt context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	req := &pb.SubscribeRequest{}

	for _, tpc := range flag.Args() {
		req.Topic = append(req.Topic, &pb.Topic{TopicName: tpc})
	}

	client := createServeiceStub(s.server)

	pubsubClient, err := client.Subscribe(ctxt, req)
	if err != nil {
		log.Fatalln("Error subscribing ", err)
	}

	for {
		msgReceived, err := pubsubClient.Recv()

		if err == io.EOF {

			continue
		}
		if err != nil {
			log.Println("err ", err)
			log.Fatalf("Failed to receive a msg %v", err)
		}

		log.Println("Message Received : ", msgReceived.Msg)

	}

	return subcommands.ExitSuccess

}
