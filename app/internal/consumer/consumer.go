package consumer

import (
	"context"
	"encoding/json"

	"test-task/internal/models"
	"test-task/internal/retry"
	"test-task/internal/service"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.Service
	retry   retry.Retrier
	log     *zap.Logger
}

func NewConsumer(cfg kafka.ReaderConfig, service *service.Service, retry retry.Retrier, log *zap.Logger) *Consumer {
	return &Consumer{
		reader:  kafka.NewReader(cfg),
		service: service,
		retry:   retry,
		log:     log,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.Info("consumer stopped by context")
				return
			}
			c.log.Error("error on reading message", zap.Error(err))
			continue
		}

		eo := new(models.ExtendedOrder)
		if err := json.Unmarshal(m.Value, eo); err != nil {
			c.log.Warn("invalid json model from message", zap.Error(err), zap.ByteString("json_model", m.Value))
			continue
		}

		if err := models.Validate(eo); err != nil {
			c.log.Warn("invalid model", zap.Error(err))
			continue
		}

		c.log.Info("creating extended order...", zap.Int64("id", eo.Order.ID))

		if err := c.retry.Do(ctx, func(attempt int) error {
			if err := c.service.CreateExtendedOrder(ctx, eo); err != nil {
				c.log.Warn("error on creating order",
					zap.Int64("id", eo.Order.ID),
					zap.Error(err),
					zap.Int("attempt", attempt),
				)
				return err
			}
			c.log.Info("order created", zap.Int("attempt", attempt), zap.Int64("id", eo.Order.ID))
			return nil
		}); err != nil {
			if ctx.Err() != nil {
				c.log.Info("consumer stopped by context")
				return
			}
			c.log.Error("failed to create order",
				zap.Int64("id", eo.Order.ID),
				zap.Error(err),
			)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
