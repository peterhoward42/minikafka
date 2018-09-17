package svr

type Server struct {
	foo int
}

func NewServer(a_number int) *Server {
	return &Server{
		foo: a_number,
	}
}
