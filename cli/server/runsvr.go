package main

import (
	"fmt"
	"time"
)

import (
	"github.com/peterhoward42/toy-kafka/protocol"
	"github.com/peterhoward42/toy-kafka/svr"
)

// This commmand launches a program that creates a Toy-Kafka server and
// mandates it to start serving.
func main() {
	// Todo: Override port from environment variables.
	// Todo: Override retentionsTime from command line argument.

	host := "localhost"
	port := protocol.DefaultPort
	svr := svr.NewServer()

	const retentionTime = time.Duration(time.Second * 5)
	svr.Serve(host, port, retentionTime)

	fmt.Println("Server Finished")
}
