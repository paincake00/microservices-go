package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/paincake00/microservices-go/payment/internal/api/payment/v1"
	"github.com/paincake00/microservices-go/payment/internal/config"
	"github.com/paincake00/microservices-go/payment/internal/service/payment"
	"github.com/paincake00/microservices-go/payment/internal/util/uuidutil"
	"github.com/paincake00/microservices-go/platform/pkg/grpc/health"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	uuidGenerator := uuidutil.NewUuidGeneratorV7()
	paymentService := payment.NewService(uuidGenerator)
	paymentHandler := v1.NewPaymentHandler(paymentService)

	srv := grpc.NewServer()

	paymentv1.RegisterPaymentServiceServer(srv, paymentHandler)

	reflection.Register(srv)

	// Регистрация эндпоинтов для healthcheck
	health.RegisterServer(srv)

	notify := make(chan error, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.Port()))
		if err != nil {
			log.Printf("Failed to listen: %v", err)

			notify <- err

			close(notify)
			return
		}

		log.Printf("Starting gRPC server on port %s", cfg.Grpc.Port())

		if err = srv.Serve(lis); err != nil {
			log.Printf("Failed to serve: %v", err)

			notify <- err
		}
		close(notify)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-interrupt:
		log.Printf("Received interrupt signal: %v", sig)
	case errNotify := <-notify:
		if errNotify != nil {
			log.Printf("Notify about suddenly error: %v", errNotify)
		}
	}

	log.Printf("Shutting down gRPC server on port %s", cfg.Grpc.Port())

	srv.GracefulStop()
}
