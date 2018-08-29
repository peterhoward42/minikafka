package client

import (
    "bufio"
    "encoding/gob"
    "errors"
    "fmt"
    "net"

    "github.com/peterhoward42/toy-kafka/svr"
)

type Producer struct {
    topic string
    rw bufio.ReadWriter
    gobEncoder *gob.Encoder
}

func NewProducer(topic string) (*Producer, error) {

    // Establish a connection to the server.
    conn, err := net.Dial("tcp", "localhost:2000")
	if err != nil {
		return nil, errors.New(
		    fmt.Sprintf("Producer failed to connect to server with error: %v",
		    err))
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	enc := gob.NewEncoder(rw)

	return &Producer{
	    topic: topic,
	    rw: *rw,
	    gobEncoder: enc,
	}, nil
}

func (prod *Producer) SendMessage(message string) (
        msg_num int, err error) {

    err = prod.gobEncoder.Encode(svr.ProduceCmd)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(
		    "Failed to encode ProduceCmd with: %v", err))
	}
	err = prod.rw.Flush()
	if err != nil {
		return -1, errors.New(fmt.Sprintf(
		    "Flush failed in SendMessage, with: %v", err))
	}

	// Read the acknowledgement before sending message payload.

	// todo - how to read and inspect the response

    return 0, fmt.Errorf("place holder")
}

