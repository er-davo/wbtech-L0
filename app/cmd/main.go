package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"test-task/internal/app"
	"test-task/internal/config"
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

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	app, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("error on creating app", zap.Error(err))
	}

	errCh := make(chan error, 1)
	go func() {
		if runErr := app.Run(ctx); runErr != nil {
			errCh <- runErr
		}
	}()

	select {
	case <-quit:
		log.Info("shutting down gracefully...")
		cancel()
	case runErr := <-errCh:
		log.Error("application exited with error", zap.Error(runErr))
		cancel()
	}

	if err := app.Shutdown(); err != nil {
		log.Error("failed to shutdown application", zap.Error(err))
	}
}
