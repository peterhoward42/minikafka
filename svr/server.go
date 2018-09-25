package svr

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations"
)

// Server *is* the toy kafka (grpc) server.
type Server struct {
	// The coupling between the server and its storage backend is governed
	// by the BackingStore interface.
	store contract.BackingStore
}

// NewServer creates and initialises a new server, but does not fire
// up the underlying grpc server.
func NewServer() *Server {
	// First implementation uses an in-memory, volatile storage
	// solution.
	return &Server{implementations.NewMemStore()}
}

// Serve mandates the server to start serving.
func (s *Server) Serve(host string, port int, retentionTime time.Duration) {
	// Bring up the gRPC server.
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterToyKafkaServer(grpcServer, s)

	// Launch a goroutine to periodically delete expired messages.
	go s.startCulling(retentionTime)

	log.Printf("Serving on port: %d", port)
	grpcServer.Serve(lis)
}

// startCulling is a run-forever function (intended to be run in its
// own goroutine), which removes messages from the backing store when their age
// exceeds *retentionTime*.
func (s *Server) startCulling(retentionTime time.Duration) {
    cullCheckFrequency := time.Duration(2) * time.Second
	ticker := time.NewTicker(cullCheckFrequency)
	for range ticker.C {
		// Note time.Add() and time.Sub() operate with differing types,
		// and the use of Add() here is deliberate. Also that you can do
		// unary-minus on the *retentionTime* time.Duration struct.
		maxAge := time.Now().Add(-retentionTime)
		// Delegate to the backing store implementation.
		err := s.store.RemoveOldMessages(maxAge); if err != nil {
            log.Fatalf("Error removing old messages: %v", err)
        }
	}
}

// gRPC HANDLER METHODS BELOW
// IF WE END UP WITH MORE THAN 2 OR 3 - SPLIT INTO SEPARATE MODULE OR EVEN
// PACKAGE

// Produce is the server's handler function for the *Produce* API call.
func (s *Server) Produce(
	ctx context.Context, req *pb.ProduceRequest) (*pb.MsgNumber, error) {

	topicStr := req.GetTopic().Topic
	messageBytes := req.GetPayload().Payload

	// Delegate storage to the backing store provider.

	msgNumber, err := s.store.Store(topicStr, messageBytes)
	if err != nil {
		log.Fatalf("Error storing message: %v", err)
	}
	return &pb.MsgNumber{MsgNumber: uint32(msgNumber)}, nil
}
