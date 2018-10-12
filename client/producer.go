package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
)

// Producer is a ToyKafkaClient client object dedicated to sending *produce*
// messages to the server.
type Producer struct {
	topic       string
    timeout     time.Duration
	clientProxy pb.ToyKafkaClient
}

// NewProducer provides a new Producer instance that is bound to a given
// host, and a given message topic.
// *host* should be of the form "myhost.com:1234".
func NewProducer(topic string, timeout time.Duration, host string) (
        *Producer, error) {

	p := &Producer{
        topic: topic,
        timeout: timeout,
    }
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	p.clientProxy = pb.NewToyKafkaClient(conn)
	return p, nil
}

// SendMessage is the primary API method for Producer, which sends
// the given message payload to the server in a Produce message.
func (p *Producer) SendMessage(messagePayload MessagePayload) (
	msgNum uint32, err error) {
	log.Printf("Producer sending msg.")
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	topic := &pb.Topic{Topic: p.topic}
	payload := &pb.Payload{Payload: messagePayload}
	produceRequest := &pb.ProduceRequest{Topic: topic, Payload: payload}
	msgNumber, err := p.clientProxy.Produce(ctx, produceRequest)
	if err != nil {
		log.Fatalf("Call to client proxy Produce() failed: %v.", err)
	}
	log.Printf("Reply message number: %v", msgNumber)
	return msgNumber.GetMsgNumber(), err
}
