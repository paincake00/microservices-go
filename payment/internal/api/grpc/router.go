package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/paincake00/microservices-go/payment/internal/api/grpc/payment/v1"
	"github.com/paincake00/microservices-go/payment/internal/service"
	"github.com/paincake00/microservices-go/platform/pkg/grpc/health"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

func NewRouter(srv *grpc.Server, paymentService service.IPaymentService) {
	paymentHandler := v1.NewPaymentHandler(paymentService)

	paymentv1.RegisterPaymentServiceServer(srv, paymentHandler)

	reflection.Register(srv)

	// Регистрация эндпоинтов для healthcheck
	health.RegisterServer(srv)
}
