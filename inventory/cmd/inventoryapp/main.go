package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	signalGo "os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/paincake00/microservices-go/inventory/internal/api/inventory/v1"
	"github.com/paincake00/microservices-go/inventory/internal/config"
	"github.com/paincake00/microservices-go/inventory/internal/repository/part/mongodb"
	partServ "github.com/paincake00/microservices-go/inventory/internal/service/part"
	"github.com/paincake00/microservices-go/inventory/pkg/mongoclient"
	"github.com/paincake00/microservices-go/platform/pkg/grpc/health"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

const (
	defaultPartsNum = 5
)

func main() {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatal("Error loading environment from config: ", errCfg)
	}

	// Получение URI для подключения к Mongo
	dbURI := cfg.Mongo.URI()
	if dbURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	// Создание подключения к mongo
	mongoClient, err := mongoclient.New(
		dbURI,
		mongoclient.ShutdownTimeout(15*time.Second),
	)
	if err != nil {
		log.Fatalf("Error connecting to mongodb: %v", err)
	}
	defer func() {
		cerr := mongoClient.Close()
		if cerr != nil {
			log.Printf("Error closing mongodb: %v", cerr)
		}
	}()
	// Получение БД
	mongoDB := mongoClient.Mongodb.Database(cfg.Mongo.DatabaseName())

	partRepository := mongodb.NewPartStorage(mongoDB)
	partService := partServ.NewService(partRepository)

	err = partService.InitParts(context.Background(), defaultPartsNum)
	if err != nil {
		log.Printf("Error initializing parts: %s", err)
		return
	}

	inventoryHandler := v1.NewInventoryHandler(partService)

	srv := grpc.NewServer()

	inventoryv1.RegisterInventoryServiceServer(srv, inventoryHandler)

	reflection.Register(srv)

	// Регистрация эндпоинтов для healthcheck
	health.RegisterServer(srv)

	// Канал для уведомления об ошибках
	notify := make(chan error, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.Port()))
		if err != nil {
			log.Printf("failed to listen: %v", err)

			notify <- err
			close(notify)
			return
		}

		log.Printf("starting gRPC server on port %s", cfg.Grpc.Port())

		if err = srv.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)

			notify <- err
		}

		close(notify)
	}()

	interrupt := make(chan os.Signal, 1)
	signalGo.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		log.Printf("Received interrupting signal: %v", sig)
	case errNotify := <-notify:
		log.Printf("Received notify error: %v", errNotify)
	}

	log.Println("Shutting down gRPC server...")

	srv.GracefulStop()
}
