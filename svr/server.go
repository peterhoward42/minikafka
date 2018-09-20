package svr

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
	"github.com/peterhoward42/toy-kafka/svr/backends"
)

// Server *is* the toy kafka (grpc) server.
type Server struct {
	store backends.BackingStore // Interface to pluggable alternatives.
}

// NewServer creates and initialises a new server, but does not fire
// up the underlying grpc server.
func NewServer() *Server {
	// First implementation uses an in-memory, volatile storage
	// solution.
	backingStore := backends.NewMemStore()
	return &Server{backingStore}
}

// Serve mandates the server to start serving.
func (s *Server) Serve(host string, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterToyKafkaServer(grpcServer, s)
	log.Printf("Serving on port: %d", port)
	grpcServer.Serve(lis)
}

// Produce is the server's implementation of one of the handler methods
// required by toykafka.proto.
func (s *Server) Produce(
	ctx context.Context, req *pb.ProduceRequest) (*pb.MsgNumber, error) {

	// Delegate storage to the backing store provider.
	topicStr := req.GetTopic().Topic
	messageBytes := req.GetPayload().Payload
	msgNumber, err := s.store.Store(topicStr, messageBytes)
	if err != nil {
		log.Fatalf("Error storing message: %v", err)
	}
	return &pb.MsgNumber{MsgNumber: msgNumber}, nil
}
