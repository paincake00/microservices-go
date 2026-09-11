package e2e

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/app"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/mongo"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/network"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/path"
)

const (
	// Данные для сети
	projectName = "e2e-inventory"

	// Параметры для контейнеров
	mongoContainer         = "e2e-mongo"
	appName                = "inventory-app"
	dockerfile             = "deploy/docker/inventory/Dockerfile"
	envFilePathInContainer = "/app/config/.env"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

var (
	envFilePath         = filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env")
	envFilePathFromRoot = filepath.Join("deploy", "compose", "inventory", ".env")
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)
	mongoAuthDB := getEnvWithLogging(ctx, testcontainers.MongoAuthDBKey)

	generatedMongo, err := mongo.NewContainer(
		ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(mongoContainer),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithUsername(mongoUsername),
		mongo.WithPassword(mongoPassword),
		mongo.WithAuthDB(mongoAuthDB),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		clearTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	// Поиск корня проекта
	projectRoot := path.GetProjectRoot()
	logger.Info(ctx, "найден корень проекта", zap.String("projectRoot", projectRoot))

	appEnv := map[string]string{
		testcontainers.AppConfigPathKey: envFilePathInContainer,
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
		testcontainers.MongoPortKey: generatedMongo.Port(),
	}

	// binds для файловой системы или volumes
	binds := []app.Mount{
		{
			Source:   filepath.Join(projectRoot, envFilePathFromRoot),
			Target:   envFilePathInContainer,
			ReadOnly: true,
		},
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(grpcPort + "/tcp").
		WithStartupTimeout(startupTimeout)

	generatedApp, err := app.NewContainer(
		ctx,
		app.WithName(appName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, dockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithBinds(binds),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		clearTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")

	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     generatedApp,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}
	return val
}
