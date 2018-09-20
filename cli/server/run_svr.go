package main

import "fmt"

import (
	"github.com/peterhoward42/toy-kafka/protocol"
	"github.com/peterhoward42/toy-kafka/svr"
)

func main() {
	// Todo: Override port from environment variables.
	host := "localhost"
	port := protocol.DefaultPort
	svr := svr.NewServer()
	svr.Serve(host, port)
	fmt.Println("Server Finished")
}
