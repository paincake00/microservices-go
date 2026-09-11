//go:build integration

package mongodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/inventory/internal/repository"
	"github.com/paincake00/microservices-go/inventory/pkg/mongoclient"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/mongo"
)

const (
	loggerLevel  = "debug"
	loggerAsJSON = true

	mongoContainer = "test-mongo"
	mongoImage     = "mongo:8.0"
	mongoUser      = "user"
	mongoPass      = "password"
	mongoDatabase  = "inventory-db"
	mongoAuthDB    = "admin"
)

type RepoSuite struct {
	suite.Suite

	container *mongo.Container

	partRepo repository.IPartRepository
}

func (s *RepoSuite) SetupSuite() {
	ctx := context.Background()

	// Инициализаия логгера
	logger.Init(loggerLevel, loggerAsJSON)

	// Создание Mongo
	container, err := mongo.NewContainer(
		ctx,
		mongo.WithNetworkName(""), // сеть не нужна
		mongo.WithContainerName(mongoContainer),
		mongo.WithImageName(mongoImage),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithUsername(mongoUser),
		mongo.WithPassword(mongoPass),
		mongo.WithAuthDB(mongoAuthDB),
		mongo.WithLogger(logger.Logger()),
	)
	s.Require().NoError(err)

	s.container = container

	// Создание подключения к mongo через клиента
	mongoClient, err := mongoclient.New(
		fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s?authSource=%s",
			mongoUser, mongoPass, s.container.Host(), s.container.Port(), mongoDatabase, mongoAuthDB,
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
	logger.Info(context.Background(), "🛑 Контейнер MongoDB остановлен")
}

func TestPartRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
