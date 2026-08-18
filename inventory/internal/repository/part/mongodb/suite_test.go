package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/paincake00/microservices-go/inventory/internal/repository"
	"github.com/paincake00/microservices-go/inventory/pkg/mongoclient"
)

const (
	mongoUser     = "user"
	mongoPass     = "password"
	mongoDatabase = "inventory-db"
)

type RepoSuite struct {
	suite.Suite

	container testcontainers.Container

	partRepo repository.IPartRepository
}

func (s *RepoSuite) SetupSuite() {
	ctx := context.Background()

	// Запрос на выкачивание образа Mongo
	req := testcontainers.ContainerRequest{
		Image:        "mongo:8.0",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": mongoUser,
			"MONGO_INITDB_ROOT_PASSWORD": mongoPass,
			"MONGO_INITDB_DATABASE":      mongoDatabase,
		},
		WaitingFor: wait.ForExec(
			[]string{
				"mongosh",
				"--quiet",
				"-u", mongoUser,
				"-p", mongoPass,
				"--authenticationDatabase", "admin",
				"--eval", "db.runCommand({ ping: 1 }).ok",
			},
		).WithStartupTimeout(60 * time.Second),
	}

	// Создание контейнера с MongoDB
	container, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	s.Require().NoError(err)

	s.container = container

	// Получение хоста
	host, err := container.Host(ctx)
	s.Require().NoError(err)

	// Получение порта
	port, err := container.MappedPort(ctx, "27017")
	s.Require().NoError(err)

	// Создание подключения к mongo через клиента
	mongoClient, err := mongoclient.New(
		fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s?authSource=admin",
			mongoUser, mongoPass, host, port.Port(), mongoDatabase,
		),
	)
	s.Require().NoError(err)

	// Полученение необходимой базы
	mongoDB := mongoClient.Mongodb.Database(mongoDatabase)

	s.partRepo = NewPartStorage(mongoDB)
}

func (s *RepoSuite) TearDownSuite() {
	err := s.container.Terminate(context.Background())
	s.Require().NoError(err)
}

func TestPartRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
