package main

import (
	"context"

	"github.com/paincake00/microservices-go/payment/internal/app"
	"github.com/paincake00/microservices-go/payment/internal/config"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(context.Background(), "failed to load config", zap.Error(err))
	}

	app.Run(cfg)

	//uuidGenerator := uuidutil.NewUuidGeneratorV7()
	//paymentService := payment.NewService(uuidGenerator)
	//paymentHandler := v1.NewPaymentHandler(paymentService)
	//
	//srv := grpc.NewServer()
	//
	//paymentv1.RegisterPaymentServiceServer(srv, paymentHandler)
	//
	//reflection.Register(srv)
	//
	//// Регистрация эндпоинтов для healthcheck
	//health.RegisterServer(srv)
	//
	//notify := make(chan error, 1)
	//
	//go func() {
	//	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.Port()))
	//	if err != nil {
	//		log.Printf("Failed to listen: %v", err)
	//
	//		notify <- err
	//
	//		close(notify)
	//		return
	//	}
	//
	//	log.Printf("Starting gRPC server on port %s", cfg.Grpc.Port())
	//
	//	if err = srv.Serve(lis); err != nil {
	//		log.Printf("Failed to serve: %v", err)
	//
	//		notify <- err
	//	}
	//	close(notify)
	//}()
	//
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	//
	//select {
	//case sig := <-interrupt:
	//	log.Printf("Received interrupt signal: %v", sig)
	//case errNotify := <-notify:
	//	if errNotify != nil {
	//		log.Printf("Notify about suddenly error: %v", errNotify)
	//	}
	//}
	//
	//log.Printf("Shutting down gRPC server on port %s", cfg.Grpc.Port())
	//
	//srv.GracefulStop()
}
