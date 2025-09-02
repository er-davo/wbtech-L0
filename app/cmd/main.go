package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"test-task/internal/app"
	"test-task/internal/config"
	"test-task/internal/database"
	"test-task/internal/logger"

	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()
	defer log.Sync()

	yamlConfigFilePath := os.Getenv("CONFIG_PATH")
	if yamlConfigFilePath == "" {
		log.Fatal("env ConfigPath is empty")
	}
	cfg, err := config.Load(yamlConfigFilePath)
	if err != nil {
		log.Fatal("error on loading config", zap.Error(err))
	}

	err = database.Migrate(cfg.App.MirgationDir, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("error on migrating database", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	app, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("error on creating app", zap.Error(err))
	}

	if err := app.Run(ctx); err != nil {
		if ctx.Err() != nil {
			log.Info("app stopped by context")
		} else {
			log.Error("app exited with error", zap.Error(err))
		}
	}
}
