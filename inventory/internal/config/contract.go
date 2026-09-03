package config

type GrpcConfig interface {
	Port() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
