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
}

func NewProducer(topic string) (*Producer, error) {

	return &Producer{
		topic:      topic
	}, nil
}

func (prod *Producer) SendMessage(message string) (
	msg_num int, err error) {

	return 0, fmt.Errorf("Not implemented yet.")
}
