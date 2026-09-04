package grpcserver

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	defaultGRPCServerAddress = ":50051"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Server struct {
	address    string
	App        *grpc.Server
	serverOpts []grpc.ServerOption
	notify     chan error
	logger     Logger
}

func NewServer(logger Logger, opts ...Option) *Server {
	srv := &Server{
		address: defaultGRPCServerAddress,
		App:     nil,
		notify:  make(chan error, 1),
		logger:  logger,
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.App = grpc.NewServer(srv.serverOpts...)

	return srv
}

func (s *Server) Start() {
	go func() {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- err

			close(s.notify)

			return
		}

		s.logger.Info(context.Background(), "Starting gRPC Server", zap.String("address", s.address))

		err = s.App.Serve(lis)
		// Проверка что буфер пуст
		if len(s.notify) == 0 {
			s.notify <- err
		}

		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	const op = "grpcserver.Shutdown"

	var shutdownErrors []error

	s.logger.Info(context.Background(), "Shutting down grpc server", zap.String("address", s.address))

	// Graceful stop gRPC Server
	s.App.GracefulStop()

	// Проверка на ошибки из notify (если был не interrupt)
	if err, ok := <-s.notify; ok {
		if err != nil && !errors.Is(err, context.Canceled) {
			s.logger.Error(context.Background(), op, zap.Error(err))
			shutdownErrors = append(shutdownErrors, err)
		}
	}

	return errors.Join(shutdownErrors...)
}
