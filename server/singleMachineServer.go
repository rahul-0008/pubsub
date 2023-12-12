package server

import (
	"context"
	"log"
	"sync"

	pb "github.com/pubsub/proto"
)

type SingleMachineServer struct {
	pb.UnimplementedPubsubServer
	m *sync.Map
}

func NewSingleMachine() pb.PubsubServer {
	return &SingleMachineServer{
		m: &sync.Map{},
	}
}

func (s *SingleMachineServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {

	subscribers, ok := s.m.Load(req.Topic.TopicName)
	log.Println("publishing the message ", subscribers)
	if ok {
		for _, subscriber := range subscribers.([]chan *pb.Message) {
			subscriber <- req.Message
		}

	} else {
		log.Println("No subscribers for the Topic ", req.Topic.TopicName)
	}

	return &pb.PublishResponse{
		Status: &pb.PublishResponse_Success_{
			Success: &pb.PublishResponse_Success{},
		},
	}, nil

}

func (s *SingleMachineServer) Subscribe(req *pb.SubscribeRequest, stream pb.Pubsub_SubscribeServer) error {
	log.Println("Subscription Received")

	subscriptionChannel := make(chan *pb.Message)
	log.Println(req.Topic)
	for _, topic := range req.Topic {
		name := topic.TopicName
		subscriptions, ok := s.m.Load(name)
		if ok {
			s.m.Store(name, append(subscriptions.([]chan *pb.Message), subscriptionChannel))
		} else {
			s.m.Store(name, []chan *pb.Message{subscriptionChannel})
		}
	}

	for {
		log.Println("New incoming message for the subscriber ")
		msg := <-subscriptionChannel
		err := stream.Send(&pb.SubscribeResponse{
			Msg: msg,
		})
		if err != nil {
			log.Println("Failed to Publish ", err)
			return err
		}

	}

}
