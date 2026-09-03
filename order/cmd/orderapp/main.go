package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	v1 "github.com/paincake00/microservices-go/order/internal/api/order/v1"
	"github.com/paincake00/microservices-go/order/internal/client/grpc/inventory"
	"github.com/paincake00/microservices-go/order/internal/client/grpc/payment"
	"github.com/paincake00/microservices-go/order/internal/config"
	orderRepo "github.com/paincake00/microservices-go/order/internal/repository/order/postgres"
	orderServ "github.com/paincake00/microservices-go/order/internal/service/order"
	"github.com/paincake00/microservices-go/order/pkg/grpcclient"
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
	"github.com/paincake00/microservices-go/platform/pkg/migrator"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

func main() {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Error loading config: %v", errCfg)
	}

	// Создаем пул соединений с Postgresql
	pg, err := pgclient.New(
		cfg.Postgres.URI(),
		pgclient.MaxOpenCons(cfg.Postgres.MaxOpenCons()),
		pgclient.MinIdleCons(cfg.Postgres.MinIdleCons()),
		pgclient.MaxConIdleTime(cfg.Postgres.MaxConIdleTime()),
		pgclient.MaxConLifetime(cfg.Postgres.MaxConLifetime()),
	)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pg.Close()

	// Создаем мигратор для управления миграциями
	migration := migrator.NewMigrator(stdlib.OpenDBFromPool(pg.Pool), cfg.Migration.MigrationDir())

	// Применяем миграцию (если уже была, то ничего не будет, Goose умный)
	err = migration.Up()
	if err != nil {
		log.Printf("failed migrations up: %v\n", err)
		return
	}

	// Создаем подключения к gRPC-серверам с функциями закрытия
	grpcInventoryAddress := net.JoinHostPort(cfg.Grpc.GetInventoryHost(), cfg.Grpc.GetInventoryPort())
	grpcPaymentAddress := net.JoinHostPort(cfg.Grpc.GetPaymentHost(), cfg.Grpc.GetPaymentPort())

	connInventory, errConnInventory := grpcclient.New(grpcInventoryAddress, cfg.Grpc.GetHealthCheckTimeout())
	if errConnInventory != nil {
		log.Printf("failed to connect to inventory: %v\n", errConnInventory)
		return
	}
	connPayment, errConnPayment := grpcclient.New(grpcPaymentAddress, cfg.Grpc.GetHealthCheckTimeout())
	if errConnPayment != nil {
		log.Printf("failed to connect to payment: %v\n", errConnPayment)
		return
	}

	// Создаем сервисы-клиенты для доступа к gRPC-методам
	inventoryService := inventory.NewService(inventoryv1.NewInventoryServiceClient(connInventory))
	paymentService := payment.NewService(paymentv1.NewPaymentServiceClient(connPayment))

	// Внедряем зависимости
	orderStorage := orderRepo.NewPostgresOrderStorage(pg)
	orderService := orderServ.NewOrderService(orderStorage, inventoryService, paymentService)
	orderHandler := v1.NewOrderHandler(orderService)

	orderMux, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Printf("Error create HTTP-mux from OpenAPI specs: %v", err)
		return
	}

	// Создаем роутер Chi
	r := chi.NewRouter()

	// Добавляем middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Http.RequestTimeout()))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderMux)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Http.ServerPort()),
		Handler:           r,
		ReadHeaderTimeout: cfg.Http.ReadHeaderTimeout(),
	}

	notify := make(chan error, 1)

	go func() {
		log.Printf("Starting HTTP Server on port %s", cfg.Http.ServerPort())

		if err = srv.ListenAndServe(); err != nil {
			log.Printf("HTTP server get suddenly error: %v", err)

			notify <- err
		}

		close(notify)
	}()

	// WAITING

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-notify:
		if err != nil {
			log.Printf("Received suddenly error from notify: %v", err)
		}
	case sig := <-interrupt:
		log.Printf("Received interrupting signal: %v", sig)
	}

	// SHUTDOWN

	ctxShutDown, cancel := context.WithTimeout(context.Background(), cfg.Http.ShutdownTimeout())
	defer cancel()

	log.Printf("Shutting down server...")

	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Printf("Shutdown HTTP server failed: %v", err)
	}
	err = connInventory.Close()
	if err != nil {
		log.Printf("Close connection inventory failed: %v", err)
	}
	err = connPayment.Close()
	if err != nil {
		log.Printf("Close connection payment failed: %v", err)
	}

	log.Printf("Server shutdown.")
}
