package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/paincake00/microservices-go/platform/pkg/logger"
)

const testsTimeout = 5 * time.Minute

var (
	suiteCtx    context.Context
	suiteCancel context.CancelFunc

	env *TestEnvironment
)

func TestInventoryE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order E2E Suite")
}

var _ = BeforeSuite(
	func() {
		logger.Init(loggerLevelValue, true)

		suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

		// Загружаем .env файл и устанавливаем переменные в окружение
		envVars, err := godotenv.Read(envFilePath)
		if err != nil {
			logger.Fatal(suiteCtx, "Не удалось загрузить .env файл", zap.Error(err))
		}

		// Устанавливаем переменные в окружение процесса
		for key, value := range envVars {
			_ = os.Setenv(key, value)
		}

		logger.Info(suiteCtx, "Запуск тестового окружения...")
		env = setupTestEnvironment(suiteCtx)
	},
)

var _ = AfterSuite(
	func() {
		logger.Info(context.Background(), "Завершение набора тестов")
		if env != nil {
			clearTestEnvironment(suiteCtx, env)
		}
		suiteCancel()
	},
)
