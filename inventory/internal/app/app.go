package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paincake00/microservices-go/inventory/internal/api/grpc"
	"github.com/paincake00/microservices-go/inventory/internal/config"
	"github.com/paincake00/microservices-go/inventory/internal/repository"
	"github.com/paincake00/microservices-go/inventory/internal/repository/part/mongodb"
	"github.com/paincake00/microservices-go/inventory/internal/service"
	partServ "github.com/paincake00/microservices-go/inventory/internal/service/part"
	"github.com/paincake00/microservices-go/inventory/pkg/grpcserver"
	"github.com/paincake00/microservices-go/inventory/pkg/mongoclient"
	"github.com/paincake00/microservices-go/platform/pkg/closer"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type useCases struct {
	partRepository repository.IPartRepository
	partService    service.IPartService
}

type servers struct {
	grpcServer *grpcserver.Server
}

func Run(cfg *config.Config) {
	cls := closer.New(logger.Logger())

	// Получение URI для подключения к Mongo
	dbURI := cfg.Mongo.URI()
	if dbURI == "" {
		logger.Fatal(context.Background(), "MONGO_URI environment variable not set")
	}

	// Создание подключения к mongo
	mongoClient, err := mongoclient.New(
		dbURI,
		mongoclient.ShutdownTimeout(15*time.Second),
	)
	if err != nil {
		logger.Fatal(context.Background(), "Error connecting to mongodb", zap.Error(err))
	}
	cls.AddNamed(
		"Mongo Client", func(ctx context.Context) error {
			return mongoClient.Close()
		},
	)
	// Получение БД
	mongoDB := mongoClient.Mongodb.Database(cfg.Mongo.DatabaseName())

	uc := initUseCases(mongoDB)
	s := initServers(cfg, uc)

	startServers(s)

	waitServersForShutdown(s, cls)
}

func initUseCases(
	mongoDB *mongo.Database,
) *useCases {
	partRepository := mongodb.NewPartStorage(mongoDB)
	partService := partServ.NewService(partRepository)

	return &useCases{
		partRepository: partRepository,
		partService:    partService,
	}
}

func initServers(cfg *config.Config, uc *useCases) *servers {
	// Создание фейковых запчастей
	err := uc.partService.InitParts(context.Background(), 5)
	if err != nil {
		logger.Fatal(context.Background(), "Error initializing parts", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(
		logger.Logger(),
		grpcserver.Addr(fmt.Sprintf(":%s", cfg.Grpc.Port())),
	)
	grpc.NewRouter(grpcServer.App, uc.partService)

	return &servers{
		grpcServer: grpcServer,
	}
}

func startServers(s *servers) {
	s.grpcServer.Start()
}

func waitServersForShutdown(s *servers, cls *closer.Closer) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		logger.Info(
			context.Background(), "Got interrupt signal for graceful shutdown", zap.String("signal", sig.String()),
		)
	case err := <-s.grpcServer.Notify():
		logger.Error(context.Background(), "app-grpcServer.Notify", zap.Error(err))
	}

	shutdownServers(s, cls)
}

func shutdownServers(s *servers, cls *closer.Closer) {
	if err := s.grpcServer.Shutdown(); err != nil {
		logger.Error(context.Background(), "app-grpcserver.Shutdown", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cls.CloseAll(ctx); err != nil {
		logger.Error(context.Background(), "app-closeAll", zap.Error(err))
	}
}
