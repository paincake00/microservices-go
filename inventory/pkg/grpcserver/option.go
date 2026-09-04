package grpcserver

type Option func(s *Server)

func Addr(address string) Option {
	return func(s *Server) {
		s.address = address
	}
}
