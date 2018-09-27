package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
)

// Consumer is a ToyKafka client object dedicated to sending *poll* messages to
// the server using gRPC.
type Consumer struct {
	topic       string
	readFrom    int               // Message number.
	clientProxy pb.ToyKafkaClient // gRPC component.
}

// NewConsumer provides a new Consumer client instance that is bound to a given
// server address, and a given message topic. The caller specifies which
// message number position they wish the subsequent polling to start.
func NewConsumer(topic string, readFrom int,
	host string, port int) (*Consumer, error) {

	p := &Consumer{topic: topic, readFrom: readFrom}
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	p.clientProxy = pb.NewToyKafkaClient(conn)
	return p, nil
}

// Poll is the primary API method for Consumer, which sends a Poll
// message to the server and returns the messages provided back to the caller.
// It also advances its internal *readFrom* position state accordingly.
func (c *Consumer) Poll() (messages []MessagePayload, err error) {
	log.Printf("Consumer Client Poll")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pollRequest := &pb.PollRequest{
		Topic: c.topic, FromMsgNumber: uint32(c.readFrom)}
	pollResponse, err := c.clientProxy.Poll(ctx, pollRequest)
	if err != nil {
		log.Fatalf("Call to client proxy Poll() failed: %v.", err)
	}
	messages = []MessagePayload{}
	payloads := pollResponse.GetPayloads()
	for _, payloadObj := range payloads {
		payload := payloadObj.GetPayload()
		messages = append(messages, payload)
	}
	c.readFrom = int(pollResponse.GetNextMsgNumber())
	log.Printf("Consumer client Poll received  %v messages", len(payloads))
	return
}
