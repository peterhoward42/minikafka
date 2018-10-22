package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pb "github.com/peterhoward42/minikafka/protocol"
)

// Producer is a MiniKafkaClient client object dedicated to sending *produce*
// messages to the server.
type Producer struct {
	topic       string
	timeout     time.Duration
	clientProxy pb.MiniKafkaClient
}

// NewProducer provides a new Producer instance that is bound to a given
// host, and a given message topic.
// *host* should be of the form "myhost.com:1234".
func NewProducer(topic string, timeout time.Duration, host string) (
	*Producer, error) {

	p := &Producer{
		topic:   topic,
		timeout: timeout,
	}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %v", err)
	}
	p.clientProxy = pb.NewMiniKafkaClient(conn)
	return p, nil
}

// SendMessage is the primary API method for Producer, which sends
// the given message payload to the server in a Produce message.
func (p *Producer) SendMessage(messagePayload MessagePayload) (
	msgNum uint32, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	topic := &pb.Topic{Topic: p.topic}
	payload := &pb.Payload{Payload: messagePayload}
	produceRequest := &pb.ProduceRequest{Topic: topic, Payload: payload}
	msgNumber, err := p.clientProxy.Produce(ctx, produceRequest)
	if err != nil {
		return 1, fmt.Errorf("client.Produce: %v", err)
	}
	return msgNumber.GetMsgNumber(), err
}
