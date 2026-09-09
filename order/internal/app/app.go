package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/paincake00/microservices-go/order/internal/api/restapi"
	clientgrpc "github.com/paincake00/microservices-go/order/internal/client/grpc"
	"github.com/paincake00/microservices-go/order/internal/client/grpc/inventory"
	"github.com/paincake00/microservices-go/order/internal/client/grpc/payment"
	"github.com/paincake00/microservices-go/order/internal/config"
	"github.com/paincake00/microservices-go/order/internal/repository"
	orderRepo "github.com/paincake00/microservices-go/order/internal/repository/order/postgres"
	"github.com/paincake00/microservices-go/order/internal/service"
	orderServ "github.com/paincake00/microservices-go/order/internal/service/order"
	"github.com/paincake00/microservices-go/order/pkg/grpcclient"
	"github.com/paincake00/microservices-go/order/pkg/httpserver"
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
	"github.com/paincake00/microservices-go/platform/pkg/closer"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"github.com/paincake00/microservices-go/platform/pkg/migrator"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

type useCases struct {
	inventoryService clientgrpc.InventoryService
	paymentService   clientgrpc.PaymentService
	orderStorage     repository.IOrderRepository
	orderService     service.IOrderService
}

type servers struct {
	httpServer *httpserver.Server
}

func Run(cfg *config.Config) {
	cls := closer.New(logger.Logger())

	// Создаем пул соединений с Postgresql
	pg, err := pgclient.New(
		cfg.Postgres.URI(),
		pgclient.MaxOpenCons(cfg.Postgres.MaxOpenCons()),
		pgclient.MinIdleCons(cfg.Postgres.MinIdleCons()),
		pgclient.MaxConIdleTime(cfg.Postgres.MaxConIdleTime()),
		pgclient.MaxConLifetime(cfg.Postgres.MaxConLifetime()),
	)
	if err != nil {
		logger.Fatal(context.Background(), "failed to connect to database", zap.Error(err))
	}
	cls.AddNamed(
		"Postgres client", func(ctx context.Context) error {
			pg.Close()
			return nil
		},
	)

	// Создаем мигратор для управления миграциями
	migration := migrator.NewMigrator(stdlib.OpenDBFromPool(pg.Pool), cfg.Migration.MigrationDir())

	// Применяем миграцию (если уже была, то ничего не будет, Goose умный)
	err = migration.Up()
	if err != nil {
		logger.Fatal(context.Background(), "failed migrations up", zap.Error(err))
	}

	// Создаем подключения к gRPC-серверам с функциями закрытия
	grpcInventoryAddress := net.JoinHostPort(cfg.Grpc.GetInventoryHost(), cfg.Grpc.GetInventoryPort())
	grpcPaymentAddress := net.JoinHostPort(cfg.Grpc.GetPaymentHost(), cfg.Grpc.GetPaymentPort())

	grpcClientInventory, errConnInventory := grpcclient.New(grpcInventoryAddress, cfg.Grpc.GetHealthCheckTimeout())
	if errConnInventory != nil {
		logger.Fatal(context.Background(), "failed to connect to inventory", zap.Error(errConnInventory))
	}
	cls.AddNamed(
		"Grpc Inventory client", func(ctx context.Context) error {
			return grpcClientInventory.Close()
		},
	)
	grpcClientPayment, errConnPayment := grpcclient.New(grpcPaymentAddress, cfg.Grpc.GetHealthCheckTimeout())
	if errConnPayment != nil {
		logger.Fatal(context.Background(), "failed to connect to payment", zap.Error(errConnPayment))
	}
	cls.AddNamed(
		"Grpc Payment client", func(ctx context.Context) error {
			return grpcClientPayment.Close()
		},
	)

	uc := initUseCases(grpcClientInventory, grpcClientPayment, pg)
	s := initServers(cfg, uc)

	startServers(s)

	waitServersForShutdown(s, cls)
}

func initUseCases(
	grpcClientInventory *grpc.ClientConn,
	grpcClientPayment *grpc.ClientConn,
	pgClient *pgclient.Client,
) *useCases {
	// Создаем сервисы-клиенты для доступа к gRPC-методам
	inventoryService := inventory.NewService(inventoryv1.NewInventoryServiceClient(grpcClientInventory))
	paymentService := payment.NewService(paymentv1.NewPaymentServiceClient(grpcClientPayment))

	// Внедряем зависимости
	orderStorage := orderRepo.NewPostgresOrderStorage(pgClient)
	orderService := orderServ.NewOrderService(orderStorage, inventoryService, paymentService)

	return &useCases{
		inventoryService: inventoryService,
		paymentService:   paymentService,
		orderStorage:     orderStorage,
		orderService:     orderService,
	}
}

func initServers(cfg *config.Config, uc *useCases) *servers {
	httpRouter, err := restapi.NewRouter(uc.orderService, cfg.Http.RequestTimeout())
	if err != nil {
		logger.Fatal(context.Background(), "failed create http router", zap.Error(err))
	}
	httpServer := httpserver.New(
		logger.Logger(),
		httpRouter,
		httpserver.Addr(fmt.Sprintf(":%s", cfg.Http.ServerPort())),
		httpserver.ReadHeaderTimeout(cfg.Http.ReadHeaderTimeout()),
		httpserver.ShutdownTimeout(cfg.Http.ShutdownTimeout()),
	)

	return &servers{
		httpServer: httpServer,
	}
}

func startServers(s *servers) {
	s.httpServer.Start()
}

func waitServersForShutdown(s *servers, cls *closer.Closer) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		logger.Info(
			context.Background(), "Got interrupt signal for graceful shutdown", zap.String("signal", sig.String()),
		)
	case err := <-s.httpServer.Notify():
		logger.Error(context.Background(), "app-httpServer.Notify", zap.Error(err))
	}

	shutdownServers(s, cls)
}

func shutdownServers(s *servers, cls *closer.Closer) {
	if err := s.httpServer.Shutdown(); err != nil {
		logger.Error(context.Background(), "app-httpserver.Shutdown", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cls.CloseAll(ctx); err != nil {
		logger.Error(context.Background(), "app-closeAll", zap.Error(err))
	}
}
