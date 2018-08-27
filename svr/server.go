package  svr

import (
	"log"
	"net"
)

type Server struct {
    foo int
}

func NewServer(a_number int) *Server {
	return &Server{
		foo: a_number,
	}
}

func (*Server) Run(msg string) {
	listener, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept() // From one tcp client.
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Accepted a connection")
		defer conn.Close()

		// Launch a per-connection command interpreter in a goroutine, that
		// consumes the incoming stream.

        interpreter := NewInterpreter()
        go interpreter.Interpret(conn)
	}
}
