package config

import "time"

type GrpcConfig interface {
	GetInventoryHost() string
	GetInventoryPort() string
	GetPaymentHost() string
	GetPaymentPort() string
	GetHealthCheckTimeout() time.Duration
}

type HttpConfig interface {
	ServerPort() string
	ReadHeaderTimeout() time.Duration
	RequestTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	DSN() string
	MaxOpenCons() int
	MinIdleCons() int
	MaxConIdleTime() time.Duration
	MaxConLifetime() time.Duration
}

type MigrationConfig interface {
	MigrationDir() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
