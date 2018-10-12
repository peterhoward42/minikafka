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

// Server *is* the toy kafka server.
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

// Serve mandates the server to start serving and also to start the automatic
// culling of expired messages.
// *host* should be of the form "myhost.com:1234".
func (s *Server) Serve(host string, retentionTime time.Duration) error {

	// Channels to receive errors back from the goroutines we launch.
	gRPCErrC := make(chan error)
	cullingErrC := make(chan error)

	// Channel to tell the message culling service to stop.
	cullingStopC := make(chan bool)

	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	go s.startGrpcServer(gRPCErrC, grpcServer, host)
	go s.startCullingService(cullingErrC, cullingStopC, retentionTime)

	// Wait forever, or, for either of the the run-forever goroutines to report
	// that it has stopped on error. When either stops on an error, explitly
	// stop the other to avoid leaking goroutines.
	var err error
	select {
	case err = <-gRPCErrC:
		cullingStopC <- true
		return fmt.Errorf("startGrpcserver: %v", err)
	case err := <-cullingErrC:
		grpcServer.GracefulStop()
		return fmt.Errorf("startCullingService: %v", err)
	}
}

//------------------------------------------------------------------------
// gRPC REQUEST HANDLERS - as per protobuf spec.
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
		return nil, fmt.Errorf("store.Store: %v", err)
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
		return nil, fmt.Errorf("store.Poll: %v", err)
	}
	payloads := []*pb.Payload{}
	for _, msg := range messages {
		payloads = append(payloads, &pb.Payload{Payload: msg})
	}
	return &pb.PollResponse{
		Payloads:    payloads,
		NewReadFrom: &pb.MsgNumber{MsgNumber: uint32(nextMsgNumber)}}, nil
}

//------------------------------------------------------------------------
// Internal helpers
//------------------------------------------------------------------------

// startGrpcServer starts listening on the requested host network interface,
// introduces the standard library gRPC server to this customer server
// wrapper, and starts it serving. If it encouters an error while running, it
// stops itself and signals the error on the error reporting channel passed in.
func (s *Server) startGrpcServer(
	errc chan<- error, grpcServer *grpc.Server, host string) {

	lis, err := net.Listen("tcp", host)
	if err != nil {
		errc <- fmt.Errorf("net.Listen: %v", err)
		return
	}
	pb.RegisterToyKafkaServer(grpcServer, s)

	err = grpcServer.Serve(lis) // Runs forever, or error encountered.
	if err != nil {
		errc <- fmt.Errorf("grpcServer.Serve: %v", err)
		return
	}
}

// startCulling periodically removes messages from the backing store when their
// age exceeds *retentionTime*. It runs forever, or, until an error occurs, or
// it receives an instruction to stop on the stop channel passed in. When
// an error occurs it signals this on the error reporting channel passed in.
func (s *Server) startCullingService(
	errc chan<- error, stopc <-chan bool, retentionTime time.Duration) {
	// If we're keeping messages until they are 50 minutes old, we check to see
	// if any have expired every 5 minutes. (one tenth of the retention time.)
	cullCheckFrequency := retentionTime / 10
	ticker := time.NewTicker(cullCheckFrequency)
	// For as long as ticks arrive...
	for range ticker.C {
		// Been instructed to stop since last tick?
		select {
		case <-stopc:
			return
		default:
		}
		// Note time.Add() and time.Sub() operate with differing types,
		// and the use of Add() here is deliberate. Also that you can do
		// unary-minus on the *retentionTime* time.Duration struct.
		maxAge := time.Now().Add(-retentionTime)
		// Delegate to the backing store implementation.
		nRemoved, err := s.store.RemoveOldMessages(maxAge)
		if err != nil {
			errc <- fmt.Errorf("store.RemoveOldMessages: %v", err)
			return
		}
		if nRemoved != 0 {
			log.Printf("Removed %d expired messages", nRemoved)
		}
	}
}
