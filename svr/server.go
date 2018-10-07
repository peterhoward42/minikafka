package svr

import (
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
// *host* should be of the form "myhost.com:1234".
// This call also starts the server's automatic removal of old messages from
// the store - based on the retention time provided.
func (s *Server) Serve(host string, retentionTime time.Duration) {
	// Bring up the gRPC server.
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterToyKafkaServer(grpcServer, s)

	// Launch a goroutine to periodically delete expired messages.
	go s.startCulling(retentionTime)

	log.Printf("Serving on host: %s", host)
	grpcServer.Serve(lis)
}

// startCulling is a run-forever function (intended to be run in its
// own goroutine), which removes messages from the backing store when their age
// exceeds *retentionTime*.
func (s *Server) startCulling(retentionTime time.Duration) {
	// If we're keeping messages until they are 50 minutes old, we check to see
	// if any have expired every 5 minutes. (one tenth of the retention time.)
	cullCheckFrequency := retentionTime / 10
	ticker := time.NewTicker(cullCheckFrequency)
	for range ticker.C {
		// Note time.Add() and time.Sub() operate with differing types,
		// and the use of Add() here is deliberate. Also that you can do
		// unary-minus on the *retentionTime* time.Duration struct.
		maxAge := time.Now().Add(-retentionTime)
		// Delegate to the backing store implementation.
		_, err := s.store.RemoveOldMessages(maxAge)
		if err != nil {
			log.Fatalf("Error removing old messages: %v", err)
		}
	}
}

//------------------------------------------------------------------------
// gRPC HANDLERS
//------------------------------------------------------------------------

// Produce is the server's handler function for the *Produce* API call.
func (s *Server) Produce(
	ctx context.Context, req *pb.ProduceRequest) (*pb.MsgNumber, error) {
	// Harvest the request details from the incoming gRPC request object,
	// then delegate the storage work to the backing store, and finally
	// package up the data to return to suit a gRPC response.

	topicStr := req.GetTopic().Topic
	messageBytes := req.GetPayload().Payload
	msgNumber, err := s.store.Store(topicStr, messageBytes)
	if err != nil {
		log.Fatalf("Error storing message: %v", err)
	}
	return &pb.MsgNumber{MsgNumber: uint32(msgNumber)}, nil
}

// Poll is the server's handler function for the *Poll* API call.
func (s *Server) Poll(ctx context.Context, req *pb.PollRequest) (
	*pb.PollResponse, error) {
	// Harvest the request details from the incoming gRPC request, then
	// delegate the retrieval of messages to the backing store, and finally,
	// package up the data to return to suit a gRPC response.

	topicStr := req.GetTopic()
	fromMsgNumber := req.GetReadFrom().GetMsgNumber()
	messages, nextMsgNumber, err := s.store.Poll(topicStr, int(fromMsgNumber))
	if err != nil {
		log.Fatalf("Error reported by backend for Poll: %v", err)
	}
	payloads := []*pb.Payload{}
	for _, msg := range messages {
		payloads = append(payloads, &pb.Payload{Payload: msg})
	}
	return &pb.PollResponse{
		Payloads:    payloads,
		NewReadFrom: &pb.MsgNumber{MsgNumber: uint32(nextMsgNumber)}}, nil
}
