package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

const grpcServerPort = 50052

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) PayOrder() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", id.String())

	return id.String(), nil
}

type IPaymentService interface {
	PayOrder() (string, error)
}

type PaymentHandler struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService IPaymentService
}

func NewPaymentHandler(paymentService IPaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (p *PaymentHandler) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (
	*paymentv1.PayOrderResponse,
	error,
) {
	id, err := p.paymentService.PayOrder()
	if err != nil {
		log.Printf("Транзакция отменена. transaction_uuid не создан: %v", err)

		return nil, status.Errorf(codes.Internal, "Транзакция отменена. transaction_uuid не создан: %v", err)
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: id,
	}, nil
}

func main() {
	paymentService := NewPaymentService()
	paymentHandler := NewPaymentHandler(paymentService)

	srv := grpc.NewServer()

	paymentv1.RegisterPaymentServiceServer(srv, paymentHandler)

	reflection.Register(srv)

	notify := make(chan error, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServerPort))
		if err != nil {
			log.Printf("Failed to listen: %v", err)

			notify <- err

			close(notify)
			return
		}

		log.Printf("Starting gRPC server on port %d", grpcServerPort)

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

	log.Printf("Shutting down gRPC server on port %d", grpcServerPort)

	srv.GracefulStop()
}
