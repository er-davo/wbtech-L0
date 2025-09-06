package service

import (
	"context"
	"test-task/internal/cache"
	"test-task/internal/models"
	"test-task/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Service struct {
	db *pgxpool.Pool

	repo repository.ExtendedOrderRepository

	cache *cache.Cache[int64, *models.ExtendedOrder]

	log *zap.Logger
}

func NewService(
	db *pgxpool.Pool,
	repo repository.ExtendedOrderRepository,
	orderCacheSize int,
	log *zap.Logger,
) *Service {
	return &Service{
		db:    db,
		repo:  repo,
		cache: cache.New[int64, *models.ExtendedOrder](orderCacheSize),
		log:   log,
	}
}

func (s *Service) LoadRecentOrdersToCache(ctx context.Context, limit int) error {
	orders, err := s.repo.GetLastExtendedOrders(ctx, limit)
	if err != nil {
		s.log.Error("failed to load last orders from db", zap.Error(err))
		return err
	}

	for _, order := range orders {
		s.cache.Add(order.ID, order)
	}

	s.log.Info("recent orders loaded to cache", zap.Int("count", len(orders)))

	return nil
}

func (s *Service) CreateExtendedOrder(ctx context.Context, eo *models.ExtendedOrder) error {
	err := s.repo.CreateExtendedOrder(ctx, eo)
	if err != nil {
		s.log.Error("failed to create order", zap.Error(err))
		return err
	}

	s.cache.Add(eo.Order.ID, eo)

	s.log.Info("order created and cached", zap.Int64("id", eo.Order.ID), zap.String("order_uid", eo.Order.OrderUID))

	return nil
}

func (s *Service) GetExtendedOrder(ctx context.Context, id int64) (*models.ExtendedOrder, error) {
	if eo, ok := s.cache.Get(id); ok {
		s.log.Info("order loaded from cache", zap.Int64("id", id))
		return eo, nil
	}

	eo, err := s.repo.GetExtendedOrder(ctx, id)
	if err != nil {
		s.log.Error("failed to load order from db", zap.Error(err), zap.Int64("id", id))
		return nil, err
	}

	s.log.Info("order loaded from db", zap.Int64("id", id))

	return eo, nil
}
