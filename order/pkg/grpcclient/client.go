package grpcclient

import (
	"fmt"
	"time"

	"github.com/paincake00/microservices-go/platform/pkg/grpc/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(address string, healthCheckTimeout time.Duration) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect to gRPC Server with address %s: %w", address, err)
	}

	// Проверка gRPC сервера на работоспособность
	err = health.Check(conn, healthCheckTimeout)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed health check: %w", err)
	}

	return conn, nil
}
