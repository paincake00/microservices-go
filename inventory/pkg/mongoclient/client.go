package mongoclient

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultOperationTimeout = 1 * time.Second // the client-side operation timeout (CSOT)
	defaultMaxPoolSize      = 100             // Maximum number of connections opened in the pool.
	defaultMinPoolSize      = 0               // Minimum number of connections opened in the pool.
	defaultMaxCons          = 2               // Maximum number of connections a pool may establish concurrently.
	defaultMaxConIdleTime   = 5 * time.Minute // The maximum number of milliseconds
	// that a connection can remain idle in the pool before being removed and closed.

	defaultShutdownTimeout = 10 * time.Second
)

type Client struct {
	Mongodb *mongo.Client

	operationTimeout time.Duration
	maxPoolSize      int
	minPoolSize      int
	maxCons          int
	maxConIdleTime   time.Duration

	shutdownTimeout time.Duration
}

func New(uri string, opts ...Option) (*Client, error) {
	mc := &Client{
		operationTimeout: defaultOperationTimeout,
		maxPoolSize:      defaultMaxPoolSize,
		minPoolSize:      defaultMinPoolSize,
		maxCons:          defaultMaxCons,
		maxConIdleTime:   defaultMaxConIdleTime,
		shutdownTimeout:  defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(mc)
	}

	if mc.maxPoolSize < 0 || mc.minPoolSize < 0 || mc.maxCons < 0 {
		return nil, fmt.Errorf("invalid options: pool size and cons nums shouldn't be less than 0")
	}

	clientOptions := options.Client().
		ApplyURI(uri).
		SetTimeout(mc.operationTimeout).
		SetMaxPoolSize(uint64(mc.maxPoolSize)).
		SetMinPoolSize(uint64(mc.minPoolSize)).
		SetMaxConnecting(uint64(mc.maxCons)).
		SetMaxConnIdleTime(mc.maxConIdleTime)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	mc.Mongodb = client

	// Context for ping timeout
	healthCheckCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping mongodb for health checking
	err = mc.Mongodb.Ping(healthCheckCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("mongoDB ping timeout or error: %w", err)
	}

	return mc, nil
}

func (m *Client) Close(parentCtx context.Context) error {
	if m.Mongodb != nil {
		ctx, cancel := context.WithTimeout(parentCtx, m.shutdownTimeout)
		defer cancel()

		return m.Mongodb.Disconnect(ctx)
	}
	return nil
}
