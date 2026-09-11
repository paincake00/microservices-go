package testcontainers

// App const
const (
	AppConfigPathKey = "CONFIG_PATH"
)

// MongoDB const
const (
	MongoImageNameKey = "MONGO_IMAGE"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)

// Postgres const
const (
	PostgresImageNameKey = "POSTGRES_IMAGE"
	PostgresHostKey      = "POSTGRES_HOST"
	PostgresPortKey      = "POSTGRES_PORT"
	PostgresDatabaseKey  = "POSTGRES_DB"
	PostgresUsernameKey  = "POSTGRES_INITDB_ROOT_USERNAME"
	PostgresPasswordKey  = "POSTGRES_INITDB_ROOT_PASSWORD" //nolint:gosec
	PostgresAuthDBKey    = "POSTGRES_AUTH_DB"
)
