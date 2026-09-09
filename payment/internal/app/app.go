package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/paincake00/microservices-go/payment/internal/api/grpc"
	"github.com/paincake00/microservices-go/payment/internal/config"
	"github.com/paincake00/microservices-go/payment/internal/service"
	"github.com/paincake00/microservices-go/payment/internal/service/payment"
	"github.com/paincake00/microservices-go/payment/internal/util/uuidutil"
	"github.com/paincake00/microservices-go/payment/pkg/grpcserver"
	"github.com/paincake00/microservices-go/platform/pkg/closer"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
)

type useCases struct {
	paymentService service.IPaymentService
}

type servers struct {
	grpcServer *grpcserver.Server
}

func Run(cfg *config.Config) {
	cls := closer.New(logger.Logger())

	uc := initUseCases()
	s := initServers(cfg, uc)

	startServers(s)

	waitServersForShutdown(s, cls)
}

func initUseCases() *useCases {
	paymentService := payment.NewService(uuidutil.NewUuidGeneratorV7())

	return &useCases{
		paymentService: paymentService,
	}
}

func initServers(cfg *config.Config, uc *useCases) *servers {
	grpcServer := grpcserver.NewServer(
		logger.Logger(),
		grpcserver.Addr(fmt.Sprintf(":%s", cfg.Grpc.Port())),
	)
	grpc.NewRouter(grpcServer.App, uc.paymentService)

	return &servers{
		grpcServer: grpcServer,
	}
}

func startServers(s *servers) {
	s.grpcServer.Start()
}

func waitServersForShutdown(s *servers, cls *closer.Closer) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		logger.Info(
			context.Background(), "Got interrupt signal for graceful shutdown", zap.String("signal", sig.String()),
		)
	case err := <-s.grpcServer.Notify():
		logger.Error(context.Background(), "app-grpcServer.Notify", zap.Error(err))
	}

	shutdownServers(s, cls)
}

func shutdownServers(s *servers, cls *closer.Closer) {
	if err := s.grpcServer.Shutdown(); err != nil {
		logger.Error(context.Background(), "app-grpcserver.Shutdown", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cls.CloseAll(ctx); err != nil {
		logger.Error(context.Background(), "app-closeAll", zap.Error(err))
	}
}
