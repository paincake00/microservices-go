package main

import (
	"fmt"
	"log"
	"net"
	"os"
	signalGo "os/signal"
	"syscall"

	partServ "github.com/paincake00/microservices-go/inventory/internal/service/part"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/paincake00/microservices-go/inventory/internal/api/inventory/v1"
	partRepo "github.com/paincake00/microservices-go/inventory/internal/repository/part"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

const (
	grpcServerPort  = 50051
	defaultPartsNum = 5
)

func main() {
	partRepository := partRepo.NewRepository()
	partService := partServ.NewService(partRepository)

	err := partService.InitParts(defaultPartsNum)
	if err != nil {
		log.Fatalf("Error initializing parts: %s", err)
	}

	inventoryHandler := v1.NewInventoryHandler(partService)

	srv := grpc.NewServer()

	inventoryv1.RegisterInventoryServiceServer(srv, inventoryHandler)

	reflection.Register(srv)

	// Канал для уведомления об ошибках
	notify := make(chan error, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServerPort))
		if err != nil {
			log.Printf("failed to listen: %v", err)

			notify <- err
			close(notify)
			return
		}

		log.Printf("starting gRPC server on port %d", grpcServerPort)

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
