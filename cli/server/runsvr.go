package main

import (
	"fmt"
	"time"
)

import (
	"github.com/peterhoward42/toy-kafka/protocol"
	"github.com/peterhoward42/toy-kafka/svr"
)

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
