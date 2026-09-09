package config

type GrpcConfig interface {
	Port() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
