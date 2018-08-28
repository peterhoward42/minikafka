package client

import (
    "bufio"
    "errors"
    "fmt"
    "net"
)

type Producer struct {
    topic string
    serverReadWriter *bufio.ReadWriter
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

	return &Producer{
	    topic: topic,
	    serverReadWriter: rw,
	}, nil
}

func (*Producer) SendMessage(message string) (msg_num int, err error) {
    return 0, fmt.Errorf("place holder")
}

