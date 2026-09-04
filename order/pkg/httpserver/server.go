package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	defaultAddress           = ":8080"
	defaultReadHeaderTimeout = 5 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Server struct {
	App               *http.Server
	address           string
	notify            chan error
	readHeaderTimeout time.Duration
	shutdownTimeout   time.Duration
	logger            Logger
}

func New(logger Logger, router http.Handler, opts ...Option) *Server {
	srv := &Server{
		App:               nil,
		address:           defaultAddress,
		notify:            make(chan error, 1),
		readHeaderTimeout: defaultReadHeaderTimeout,
		shutdownTimeout:   defaultShutdownTimeout,
		logger:            logger,
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.App = &http.Server{
		Addr:              srv.address,
		Handler:           router,
		ReadHeaderTimeout: srv.readHeaderTimeout,
	}

	return srv
}

func (s *Server) Start() {
	go func() {
		s.logger.Info(context.Background(), "Starting http server", zap.String("address", s.address))
		err := s.App.ListenAndServe()

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
	const op = "httpserver.Shutdown"

	var shutdownErrors []error

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	s.logger.Info(context.Background(), "Shutting down http server", zap.String("address", s.address))

	// Завершение сервера и отлов ошибки при завершении
	err := s.App.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(context.Background(), op, zap.Error(err))
		shutdownErrors = append(shutdownErrors, err)
	}

	// Проверка на ошибки из notify (если был не interrupt)
	if err, ok := <-s.notify; ok {
		if err != nil && !errors.Is(err, context.Canceled) {
			shutdownErrors = append(shutdownErrors, err)
		}
	}

	return errors.Join(shutdownErrors...)
}
