package main

import (
	"log"

	"github.com/paincake00/microservices-go/order/internal/app"
	"github.com/paincake00/microservices-go/order/internal/config"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
)

func main() {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Error loading environment from config: %v", errCfg)
	}

	logger.Init(cfg.Logger.Level(), cfg.Logger.AsJson())
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("Error syncing logger: %v", err)
		}
	}()

	app.Run(cfg)
}
