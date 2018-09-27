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
// the server.
type Consumer struct {
	topic             string
	fromMessageNumber int
	clientProxy       pb.ToyKafkaClient
}

// NewConsumer provides a new Consumer instance that is bound to a given
// server address, and a given message topic.
func NewConsumer(topic string, fromMessageNumber int,
	host string, port int) (*Consumer, error) {

	p := &Consumer{topic: topic, fromMessageNumber: fromMessageNumber}
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
// message citing the message number from which to read.
func (c *Consumer) Poll() (messages []MessagePayload, err error) {
	log.Printf("Consumer sending a Poll msg.")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pollRequest := &pb.PollRequest{
		Topic: c.topic, FromMsgNumber: uint32(c.fromMessageNumber)}
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
	c.fromMessageNumber = int(pollResponse.GetNextMsgNumber())
	log.Printf("Received %v messages", len(payloads))
	return
}
