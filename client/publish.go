package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/subcommands"
	pb "github.com/pubsub/proto"
)

type PublishCmd struct {
	server string
	topic  pb.Topic
}

func (p *PublishCmd) Name() string {
	return "publish"
}
func (p *PublishCmd) Synopsis() string {
	return "Publsih a command to a topic mentioned"
}

func (p *PublishCmd) Usage() string {
	return `publish --topic=<topic> <message>

	Publish message to the topic
  
  `

}
func (p *PublishCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.topic.TopicName, "topic", "", "Topic to publish a message")
}

func (p *PublishCmd) Execute(ctxt context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	request := &pb.PublishRequest{}
	request.Topic = &p.topic
	request.Message = &pb.Message{Content: f.Args()[0]}

	client := createServeiceStub(p.server)
	res, err := client.Publish(ctxt, request)
	if err != nil {
		log.Fatalln("Error while publising ", err)
	}
	switch res.Status.(type) {
	case *pb.PublishResponse_Failure_:
		log.Fatalln("Publish message error ", res.GetStatus())
	}

	return subcommands.ExitSuccess

}
