package main

import (
	"log"
	"os"
	"time"

	"github.com/peterhoward42/toy-kafka/svr"
)

// This commmand-line program instantiates a Toy-Kafka server and
// mandates it to start serving.
func main() {

	host, retentionTime := readEnvironmentVariables()

	svr := svr.NewServer()
	// Server forever, or until an error condition.
	log.Printf("Launching server on: %v", host)
	err := svr.Serve(host, retentionTime)
	if err != nil {
		log.Fatalf("svr.Serve: %s", err)
	}

	log.Print("Server Finished")
}

// readEnvironmentVariables fetches the configuration parameters required
// to run the server from environment variables.
// It treats their absence as a fatal error.
func readEnvironmentVariables() (host string, retentionTime time.Duration) {

	const hostEnvVar string = "TOYKAFKA_HOST"
	const retentionEnvVar string = "TOYKAFKA_RETENTIONTIME"

	host = os.Getenv(hostEnvVar)
	rt := os.Getenv(retentionEnvVar)

	if host == "" {
		log.Fatalf("Please set the %s environment variable\n"+
			"E.g. :9999", hostEnvVar)
	}
	if rt == "" {
		log.Fatalf("Please set the %s environment variable\n"+
			"E.g. 3s or 20m", retentionEnvVar)
	}

	retentionTime, err := time.ParseDuration(rt)
	if err != nil {
		log.Fatalf("Error parsing this retention time (%s) from \n"+
			"the %s environment variable: %s", rt, retentionEnvVar, err)
	}
	return host, retentionTime
}
