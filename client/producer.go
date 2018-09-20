package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
)

// Producer provides a convenience API to a ToyKafkaClient to simplify the
// process of sending *Produce* messages.
type Producer struct {
	topic       string
	clientProxy pb.ToyKafkaClient
}

// NewProducer provides a new Producer instance that is bound to a given
// server address, and a given message topic.
func NewProducer(topic string, host string, port int) (*Producer, error) {

	p := &Producer{topic: topic}
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	p.clientProxy = pb.NewToyKafkaClient(conn)
	return p, nil
}

// SendMessage is the primary API method for Producer, which sends a Produce
// message comprising the given string to the server.
func (p *Producer) SendMessage(message string) (msgNum int32, err error) {
	log.Printf("Sending msg.")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payloadBytes := []byte(message)
	topic := &pb.Topic{Topic: p.topic}
	payload := &pb.Payload{Payload: payloadBytes}
	produceRequest := &pb.ProduceRequest{Topic: topic, Payload: payload}
	msgNumber, err := p.clientProxy.Produce(ctx, produceRequest)
	if err != nil {
		log.Fatalf("Call to client proxy Produce() failed: %v.", err)
	}
	log.Printf("Message acknowledged.")
	return msgNumber.GetMsgNumber(), err
}
