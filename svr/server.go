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
		go func(connection net.Conn) {
            // Do something on stuff that arrives here.
			connection.Close()
		}(conn)
	}
}
