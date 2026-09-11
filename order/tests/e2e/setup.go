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
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/network"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/path"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/postgres"
)

const (
	// Данные для сети
	projectName = "e2e-order"

	// Параметры для контейнеров
	postgresContainer      = "e2e-postgres"
	appName                = "order-app"
	dockerfile             = "deploy/docker/order/Dockerfile"
	envFilePathInContainer = "/app/config/.env"

	// Переменные окружения приложения
	httpPortKey = "HTTP_SERVER_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

var (
	envFilePath         = filepath.Join("..", "..", "..", "deploy", "compose", "order", ".env")
	envFilePathFromRoot = filepath.Join("deploy", "compose", "order", ".env")
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network  *network.Network
	Postgres *postgres.Container
	App      *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	postgresUsername := getEnvWithLogging(ctx, testcontainers.PostgresUsernameKey)
	postgresPassword := getEnvWithLogging(ctx, testcontainers.PostgresPasswordKey)
	postgresImageName := getEnvWithLogging(ctx, testcontainers.PostgresImageNameKey)
	postgresDatabase := getEnvWithLogging(ctx, testcontainers.PostgresDatabaseKey)

	generatedPostgres, err := postgres.NewContainer(
		ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(postgresContainer),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithUsername(postgresUsername),
		postgres.WithPassword(postgresPassword),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		clearTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер Postgres", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Postgres успешно запущен")

	httpPort := getEnvWithLogging(ctx, httpPortKey)

	projectRoot := path.GetProjectRoot()
	logger.Info(ctx, "найден корень проекта", zap.String("projectRoot", projectRoot))

	appEnv := map[string]string{
		testcontainers.AppConfigPathKey: envFilePathInContainer,
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.PostgresHostKey: generatedPostgres.Config().ContainerName,
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
	waitStrategy := wait.ForListeningPort(httpPort + "/tcp").
		WithStartupTimeout(startupTimeout)

	generatedOrder, err := app.NewContainer(
		ctx,
		app.WithName(appName),
		app.WithPort(httpPort),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithDockerfile(projectRoot, dockerfile),
		app.WithEnv(appEnv),
		app.WithBinds(binds),
		// app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		clearTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")

	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		App:      generatedOrder,
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
