package e2e

import (
	"context"

	"go.uber.org/zap"

	"github.com/paincake00/microservices-go/platform/pkg/logger"
)

func clearTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		err := env.App.Terminate(ctx)
		if err != nil {
			logger.Error(ctx, "не удалось остановить app контейнер", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер приложения остановлен")
		}
	}

	if env.Mongo != nil {
		err := env.Mongo.Terminate(ctx)
		if err != nil {
			logger.Info(ctx, "не удалось остановить контейнер MongoDB", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер MongoDB остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Сеть удалена")
		}
	}
}
