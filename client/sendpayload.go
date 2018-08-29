package client

import (
    "github.com/peterhoward42/toy-kafka/svr"
)

type SendPayload struct {
    Command svr.CommandCode
    Message string
}

func NewSendPayload(message string) *SendPayload {
    return &SendPayload{
        Command: svr.ProduceCmd,
        Message: message,
    }
}