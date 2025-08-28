package app

import (
	"context"
	"fmt"
	"net/http"
	"test-task/internal/config"
	"test-task/internal/consumer"
	"test-task/internal/database"
	"test-task/internal/handler"
	"test-task/internal/repository"
	"test-task/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type App struct {
	log *zap.Logger
	cfg *config.Config

	db *pgxpool.Pool

	consumer *consumer.Consumer
	server   *echo.Echo
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())

	e.Static("/", "public")

	retrier := newServiceRetrier(cfg.Retry, isRetryableFunc)

	repo := repository.NewExtendedOrderRepository(db)
	service := service.NewService(
		db,
		repo,
		cfg.Service.CacheSize,
		log,
	)

	if err := service.LoadRecentOrdersToCache(ctx, cfg.Service.CacheSize); err != nil {
		return nil, fmt.Errorf("failed to load recent orders to cache: %w", err)
	}

	handler := handler.NewHandler(service, retrier, log)
	handler.RegisterRoutes(e)
	consumer := consumer.NewConsumer(kafka.ReaderConfig{
		Topic:   cfg.Kafka.Topic,
		Brokers: cfg.Kafka.Brokers,
	}, service, retrier, log)

	return &App{
		cfg:      cfg,
		log:      log,
		db:       db,
		consumer: consumer,
		server:   e,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go a.consumer.Run(ctx)

	if err := a.server.Start(":" + a.cfg.App.Port); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Shutdown() error {
	if err := a.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer cancelTimeout()
	if err := a.server.Shutdown(ctxTimeout); err != nil {
		return fmt.Errorf("failed to shutdown echo server: %w", err)
	}

	a.db.Close()

	return nil
}
