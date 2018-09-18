package svr

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/peterhoward42/toy-kafka/protocol"
)

// Server *is* the toy kafka (grpc) server.
type Server struct {
}

// NewServer creates and initialises a new server, but does not fire
// up the underlying grpc server.
func NewServer() *Server {
	return &Server{}
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
	// Placeholder code - returns a hard-coded result.
	return &pb.MsgNumber{MsgNumber: 42}, nil
}
