package pgclient

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// max opened cons (всего может быть 30 подключений: все могут открыты)
	defaultMaxOpenCons = 30
	// min idle cons (может быть минимум ждущих 5 подключений)
	defaultMinIdleCons = 5
	// idle conn lifetime (ждущее живет столько, а потом закрывается)
	defaultMaxConIdleTime = 5 * time.Second
	// conn lifetime (столько живет открытое соединение, а потом закрывается)
	defaultMaxConLifetime = 30 * time.Second
)

type PostgresClient struct {
	maxOpenCons    int
	minIdleCons    int
	maxConIdleTime time.Duration
	maxConLifetime time.Duration

	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func New(url string, opts ...Option) (*PostgresClient, error) {
	pg := &PostgresClient{
		maxOpenCons:    defaultMaxOpenCons,
		minIdleCons:    defaultMinIdleCons,
		maxConIdleTime: defaultMaxConIdleTime,
		maxConLifetime: defaultMaxConLifetime,
	}

	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres url: %w", err)
	}

	maxCons, err := safeIntToInt32(pg.maxOpenCons)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres max open cons: %w", err)
	}
	poolConfig.MaxConns = maxCons

	minIdleCons, err := safeIntToInt32(pg.minIdleCons)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres min idle cons: %w", err)
	}
	poolConfig.MinIdleConns = minIdleCons

	poolConfig.MaxConnIdleTime = pg.maxConIdleTime
	poolConfig.MaxConnLifetime = pg.maxConLifetime

	pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	// Context for ping timeout
	healthCheckCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping to postgres
	if err = pg.Pool.Ping(healthCheckCtx); err != nil {
		return nil, fmt.Errorf("postgres ping timeout or error: %w", err)
	}

	return pg, nil
}

func (p *PostgresClient) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func safeIntToInt32(v int) (int32, error) {
	if v < math.MaxInt32 || v > math.MaxInt32 {
		return 0, fmt.Errorf("value %d overflows int32", v)
	}

	return int32(v), nil
}
