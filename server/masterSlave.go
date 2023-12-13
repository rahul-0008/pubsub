package server

import (
	"io"
	"log"
	"sync"

	"context"

	pb "github.com/pubsub/proto"
	"google.golang.org/grpc"
)

const slave = "SlaveTag"

type MasterServer struct {
	pb.UnimplementedMasterServiceServer
	salves *sync.Map
	SingleMachineServer
}

func (m *MasterServer) SubscribeMaster(req *pb.SlaveSubscribeRequest, stream pb.MasterService_SubscribeMasterServer) error {
	subscriber := make(chan *pb.SlaveSubscribeResponse)

	slvs, ok := m.salves.Load(slave)
	if ok {
		m.salves.Store(slave, append(slvs.([]chan *pb.SlaveSubscribeResponse), subscriber))
	} else {
		m.salves.Store(slave, make([]chan *pb.SlaveSubscribeResponse, 0))
	}

	for {
		msg := <-subscriber
		err := stream.Send(&pb.SlaveSubscribeResponse{
			Topic:   msg.Topic,
			Message: msg.Message,
		})
		if err != nil {
			log.Fatalln("Send to the subscriber failed")
			return err
		}

	}

}

func (m *MasterServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	salveSubscriptions, _ := m.salves.Load(slave)

	log.Println(salveSubscriptions)

	for _, slv := range salveSubscriptions.([]chan *pb.SlaveSubscribeResponse) {
		slv <- &pb.SlaveSubscribeResponse{
			Topic:   req.Topic,
			Message: req.Message,
		}
	}

	// also send to the subscribers of a topic in master server
	return m.SingleMachineServer.Publish(ctx, req)
}

func NewMasterServer() *MasterServer {
	slvs := &sync.Map{}
	slvs.Store(slave, []chan *pb.SlaveSubscribeResponse{})
	return &MasterServer{
		SingleMachineServer: SingleMachineServer{
			m: &sync.Map{},
		},
		salves: slvs,
	}
}

// Slave Server

type SlaveServer struct {
	MasterServer
	masterService  pb.MasterServiceClient
	pubsubServeice pb.PubsubClient
}

func NewSlaveServer(conn *grpc.ClientConn) *SlaveServer {
	slv := &SlaveServer{
		MasterServer:   *NewMasterServer(),
		masterService:  pb.NewMasterServiceClient(conn),
		pubsubServeice: pb.NewPubsubClient(conn),
	}

	err := slv.init()
	if err != nil {
		log.Fatalln("Could not initialise Slave server")
	}

	return slv
}

func (s *SlaveServer) init() error {
	stream, err := s.masterService.SubscribeMaster(context.Background(), &pb.SlaveSubscribeRequest{})
	if err != nil {
		return err
	}

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				continue
			}
			if err != nil {
				log.Fatalln(err)
			}
			subscriptions, _ := s.salves.Load(slave)
			for _, subscriber := range subscriptions.([]chan *pb.SlaveSubscribeResponse) {
				subscriber <- res
			}

			chs, ok := s.m.Load(res.Topic.TopicName)
			if ok {
				for _, ch := range chs.([]chan *pb.Message) {
					ch <- res.Message
				}
			}

		}
	}()

	return nil
}

func (s *SlaveServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	log.Println("Publish request on Slave recieved ReRouting to Master")
	return s.pubsubServeice.Publish(ctx, req)

}
