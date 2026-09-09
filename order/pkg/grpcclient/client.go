package grpcclient

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/paincake00/microservices-go/platform/pkg/grpc/health"
)

func New(address string, healthCheckTimeout time.Duration) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect to gRPC Server with address %s: %w", address, err)
	}

	// Проверка gRPC сервера на работоспособность (Health check)
	if errHc := health.Check(conn, healthCheckTimeout); errHc != nil {
		errCls := conn.Close() // закрытия подключения при ошибке (истек таймаут)
		return nil, fmt.Errorf("failed health check: %w", errors.Join(errHc, errCls))
	}

	return conn, nil
}
