package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/paincake00/microservices-go/order/internal/migrator"
	orderRepo "github.com/paincake00/microservices-go/order/internal/repository/order/memory"
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/paincake00/microservices-go/order/internal/api/order/v1"
	"github.com/paincake00/microservices-go/order/internal/client/grpc/inventory"
	"github.com/paincake00/microservices-go/order/internal/client/grpc/payment"
	orderServ "github.com/paincake00/microservices-go/order/internal/service/order"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

const (
	grpcInventoryAddress = "localhost:50051"
	grpcPaymentAddress   = "localhost:50052"

	httpPort          = 8080
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	maxOpenCons    = 30
	minIdleCons    = 5
	maxConIdleTime = 5 * time.Second
	maxConLifetime = 30 * time.Second
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	// Создаем пул соединений с Postgresql
	pg, err := pgclient.New(
		os.Getenv("DB_URI"),
		pgclient.MaxOpenCons(maxOpenCons),
		pgclient.MinIdleCons(minIdleCons),
		pgclient.MaxConIdleTime(maxConIdleTime),
		pgclient.MaxConLifetime(maxConLifetime),
	)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pg.Close()

	// Создаем мигратор для управления миграциями
	migration := migrator.NewMigrator(stdlib.OpenDBFromPool(pg.Pool), os.Getenv("MIG_DIR"))

	// Применяем миграцию (если уже была, то ничего не будет, Goose умный)
	err = migration.Up()
	if err != nil {
		log.Printf("failed migrations up: %v\n", err)
		return
	}

	// Создаем подключения к gRPC-серверам с функциями закрытия
	connInventory, closeConnInventory := createNewGrpcClient(grpcInventoryAddress)
	connPayment, closeConnPayment := createNewGrpcClient(grpcPaymentAddress)

	// Создаем сервисы-клиенты для доступа к gRPC-методам
	inventoryService := inventory.NewService(inventoryv1.NewInventoryServiceClient(connInventory))
	paymentService := payment.NewService(paymentv1.NewPaymentServiceClient(connPayment))

	// Внедряем зависимости
	orderStorage := orderRepo.NewInMemoryOrderStorage()
	orderService := orderServ.NewOrderService(orderStorage, inventoryService, paymentService)
	orderHandler := v1.NewOrderHandler(orderService)

	orderMux, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Error create HTTP-mux from OpenAPI specs: %v", err)
	}

	// Создаем роутер Chi
	r := chi.NewRouter()

	// Добавляем middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderMux)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	notify := make(chan error, 1)

	go func() {
		log.Printf("Starting HTTP Server on port %d", httpPort)

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

	ctxShutDown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Printf("Shutting down server...")

	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Printf("Shutdown HTTP server failed: %v", err)
	}
	err = closeConnInventory()
	if err != nil {
		log.Printf("Close connection inventory failed: %v", err)
	}
	err = closeConnPayment()
	if err != nil {
		log.Printf("Close connection payment failed: %v", err)
	}

	log.Printf("Server shutdown.")
}

func createNewGrpcClient(address string) (*grpc.ClientConn, func() error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to gRPC Server with address %s: %v", address, err)
	}
	return conn, func() error { return conn.Close() }
}
