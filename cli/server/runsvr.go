package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/peterhoward42/toy-kafka/cli"
	"github.com/peterhoward42/toy-kafka/svr"
)

const defaultRetentionTime = "5m"

// This commmand-line program instantiates a Toy-Kafka server and
// mandates it to start serving.
func main() {

	// Harvest some configuration from environment variables.
	var host string
	if host = os.Getenv("TOYKAFKA_HOST"); host == "" {
		log.Printf("TOYKAFKA_HOST environment variable is not set, "+
			"so using default: %s", cli.DefaultHost)
		host = cli.DefaultHost
	}

	var retentionTime time.Duration
	rt := os.Getenv("TOYKAFKA_RETENTIONTIME")
	if rt == "" {
		log.Printf("TOYKAFKA_RETENTIONTIME environment variable is not set, "+
			"so using default: %s", defaultRetentionTime)
		retentionTime, _ = time.ParseDuration(defaultRetentionTime)
	} else {
		var err error
		retentionTime, err = time.ParseDuration(rt)
		if err != nil {
			log.Fatalf("Error parsing retention time: %s\n"+
				"should be something like 5m or 2d, or 1.3s", rt)
		}
	}

	svr := svr.NewServer()
	svr.Serve(host, retentionTime)

	fmt.Println("Server Finished")
}
