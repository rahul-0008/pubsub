package main

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/pubsub/proto"
	"google.golang.org/grpc"
)

var (
	words = strings.Split("This is a test akhjgsakd sk.djcb.ksadj", "")
	topic = &pb.Topic{
		TopicName: "TestTopic",
	}

	step  = flag.Int("step", 50, "Number of clients to add in each batch")
	lower = flag.Int("lower", 0, "Lower bound of numbe rof clients to be handeled")
	upper = flag.Int("upper", 2000, "Upper bound of numbe rof clients to be handeled")

	timeout = flag.String("timeout", "10s", "Secs to timeput after inactivty")
	output  = flag.String("output", "results.csv", "path to output the result")
)

type clientResult struct {
	incorrect bool
}

type batchResult struct {
	size      int
	succeeded int
	timeout   int
	incorrect int
	duration  time.Duration
}

func launchClient(client pb.PubsubClient, r chan *clientResult) {
	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{Topic: []*pb.Topic{topic}})
	if err != nil {
		log.Fatalln("cannot connect to server ", err)
	}

	for {
		c := &clientResult{}

		for _, w := range words {
			res, err := stream.Recv()

			if err != nil {
				return
			}

			if res.Msg.Content != w {
				log.Println("Incorrect Message : got ", res.Msg.Content, " want ", w)
				c.incorrect = true
				break
			}
		}
		r <- c
	}

}

func writeBatchResult(b batchResult, o *csv.Writer) {
	log.Println("Number of clients ", b.size)
	log.Println("Nuber of suceeded ", b.succeeded)
	log.Println("Nuber of Incorrect ", b.incorrect)
	log.Println("Number of timeout ", b.timeout)
	log.Println("Total time elapsed ", b.timeout)

	err := o.Write([]string{
		strconv.Itoa(b.size),
		strconv.Itoa(b.succeeded),
		strconv.Itoa(b.timeout),
		b.duration.String(),
	})
	if err != nil {
		log.Fatalln("cannot write to csv ", err)
	}
	o.Flush()
	if o.Error() != nil {
		log.Fatalln("cannot write to csv ", o.Error())
	}
}
func main() {

	flag.Parse()
	svr := flag.Args() // server's  Ip will be parsed
	log.Println(svr)
	connection := make([]*grpc.ClientConn, len(svr))
	for i, s := range svr {
		conn, err := grpc.Dial(s, grpc.WithTimeout(10*time.Second), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalln("Client connection was not established with the server ", err)
		}
		defer conn.Close()
		connection[i] = conn // save the connection for that particualr server

	}
	// now make pubsub client with with the grpc client connection

	pub := make([]pb.PubsubClient, len(connection))
	sub := make([]pb.PubsubClient, len(connection))

	for i, conn := range connection {
		pub[i] = pb.NewPubsubClient(conn)
		sub[i] = pb.NewPubsubClient(conn)
	}

	// create the output file
	f, err := os.Create(*output)
	if err != nil {
		log.Fatalln("Erroe creating the oputputfile", err)
	}
	defer f.Close()

	// create the writer to the file
	o := csv.NewWriter(f)

	// write out the header for the csv file
	o.Write([]string{
		"size",
		"succeeded",
		"timeout",
		"duration",
	})

	r := make(chan *clientResult)

	currentLoad := *lower

	// spin up the clients in batch of *step size
	for currentLoad < *upper {
		for i := currentLoad; i < currentLoad+*step; i++ {
			go launchClient(sub[i%len(sub)], r)
		}

		currentLoad += *step
		time.Sleep(time.Second * 2)

		// publish the messages

		for i, wrd := range words {
			_, err := pub[i%len(pub)].Publish(context.Background(), &pb.PublishRequest{
				Topic: topic,
				Message: &pb.Message{
					Content: wrd,
				},
			})

			if err != nil {
				log.Fatalln("Could not publis message ", wrd, " ", err)
			}
		}

		b := batchResult{}
		b.size = currentLoad
		start := time.Now()
		to, err := time.ParseDuration(*timeout)
		if err != nil {
			log.Fatalln("cannot parse timeout ", err)
		}
		for i := 0; i < currentLoad; i++ {
			select {
			case c := <-r:
				b.succeeded++
				if c.incorrect {
					b.incorrect++
				}
			case <-time.After(to):
				b.timeout++
				start = start.Add(to)
			}
		}
		end := time.Now()
		b.duration = end.Sub(start)
		writeBatchResult(b, o)

		time.Sleep(time.Second * 2)
	}

}
