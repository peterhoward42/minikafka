package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
)

// Consumer is a ToyKafka client object dedicated to sending *poll* messages to
// the server using gRPC.
type Consumer struct {
	topic       string
	readFrom    int // Message number.
	timeout     time.Duration
	clientProxy pb.ToyKafkaClient // gRPC component.
}

// NewConsumer provides a new Consumer client instance that is bound to a given
// host, and a given message topic. The caller specifies which
// message number read-from position they wish the subsequent polling to start.
// *host* should be of the form "myhost.com:1234".
func NewConsumer(topic string, readFrom int, timeout time.Duration,
	host string) (*Consumer, error) {

	p := &Consumer{topic: topic, readFrom: readFrom, timeout: timeout}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %v", err)
	}
	p.clientProxy = pb.NewToyKafkaClient(conn)
	return p, nil
}

// Poll is the primary API method for Consumer, which sends a Poll
// message to the server and returns the messages provided back to the caller.
// It also advances its internal *readFrom* position state accordingly (ready
// for the next poll), and additionally notifies the caller of this this in its
// return values.
func (c *Consumer) Poll() (
	messages []MessagePayload, newReadFrom int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Derive the message number to read from the state held in this consumer
	// object.
	readFrom := &pb.MsgNumber{MsgNumber: uint32(c.readFrom)}

	pollRequest := &pb.PollRequest{
		Topic: c.topic, ReadFrom: readFrom}
	pollResponse, err := c.clientProxy.Poll(ctx, pollRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("client.Poll: %v", err)
	}

	// Capture the messages to return.
	messages = []MessagePayload{}
	payloads := pollResponse.GetPayloads()
	for _, payloadObj := range payloads {
		payload := payloadObj.GetPayload()
		messages = append(messages, payload)
	}

	// Update the newReadFrom message number, ready for the next poll.
	c.readFrom = int(pollResponse.GetNewReadFrom().GetMsgNumber())
	newReadFrom = c.readFrom
	return
}
