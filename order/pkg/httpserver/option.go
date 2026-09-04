package httpserver

import "time"

type Option func(s *Server)

func Addr(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

func ReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readHeaderTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
