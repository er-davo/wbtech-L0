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

	ordersRepo   repository.OrdersRepository
	itemsRepo    repository.ItemsRepository
	paymentRepo  repository.PaymentRepository
	deliveryRepo repository.DeliveryRepository

	cache *cache.Cache[int, *models.ExtendedOrder]

	log *zap.Logger
}

func NewService(
	db *pgxpool.Pool,
	oerdersRepo repository.OrdersRepository,
	itemsRepo repository.ItemsRepository,
	paymentRepo repository.PaymentRepository,
	deliveryRepo repository.DeliveryRepository,
	orderCacheSize int,
	log *zap.Logger,
) *Service {
	return &Service{
		db:           db,
		ordersRepo:   oerdersRepo,
		itemsRepo:    itemsRepo,
		paymentRepo:  paymentRepo,
		deliveryRepo: deliveryRepo,
		cache:        cache.New[int, *models.ExtendedOrder](orderCacheSize),
		log:          log,
	}
}

func (s *Service) LoadRecentOrdersToCache(ctx context.Context, limit int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			s.cache.Clear()
			tx.Rollback(ctx)
		}
	}()

	orders, err := s.ordersRepo.GetLastN(ctx, tx, limit)
	if err != nil {
		s.log.Error("failed to get last orders", zap.Error(err))
		return err
	}

	if len(orders) == 0 {
		s.log.Info("no recent orders to cache")
		return nil
	}

	orderIDs := make([]int, 0, len(orders))
	for _, o := range orders {
		orderIDs = append(orderIDs, o.ID)
	}

	deliveries, err := s.deliveryRepo.GetByOrderIDs(ctx, tx, orderIDs)
	if err != nil {
		return err
	}
	deliveryMap := make(map[int]*models.Delivery, len(deliveries))
	for _, d := range deliveries {
		deliveryMap[d.ID] = d
	}

	payments, err := s.paymentRepo.GetByOrderIDs(ctx, tx, orderIDs)
	if err != nil {
		return err
	}
	paymentMap := make(map[int]*models.Payment, len(payments))
	for _, p := range payments {
		paymentMap[p.ID] = p
	}

	items, err := s.itemsRepo.GetByOrderIDs(ctx, tx, orderIDs)
	if err != nil {
		return err
	}
	itemsMap := make(map[int][]*models.Item)
	for _, item := range items {
		itemsMap[item.OrderID] = append(itemsMap[item.OrderID], item)
	}

	for _, o := range orders {
		eo := &models.ExtendedOrder{
			Order:    *o,
			Delivery: *deliveryMap[o.DeliveryID],
			Payment:  *paymentMap[o.PaymentID],
			Items:    itemsMap[o.ID],
		}
		s.cache.Add(o.ID, eo)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	s.log.Info("recent orders loaded to cache", zap.Int("count", len(orders)))

	return nil
}

func (s *Service) CreateExtendedOrder(ctx context.Context, eo *models.ExtendedOrder) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = s.deliveryRepo.Create(ctx, tx, &eo.Delivery)
	if err != nil {
		s.log.Error("failed to create delivery", zap.Error(err))
		return err
	}

	err = s.paymentRepo.Create(ctx, tx, &eo.Payment)
	if err != nil {
		s.log.Error("failed to create payment", zap.Error(err))
		return err
	}

	eo.DeliveryID = eo.Delivery.ID
	eo.PaymentID = eo.Payment.ID

	err = s.ordersRepo.Create(ctx, tx, &eo.Order)
	if err != nil {
		s.log.Error("failed to create order", zap.Error(err))
		return err
	}

	for _, item := range eo.Items {
		item.OrderID = eo.Order.ID
	}

	err = s.itemsRepo.CreateItems(ctx, tx, eo.Items)
	if err != nil {
		s.log.Error("failed to create items", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	s.cache.Add(eo.Order.ID, eo)

	s.log.Info("order created and cached", zap.Int("id", eo.Order.ID), zap.String("order_uid", eo.Order.OrderUID))

	return nil
}

func (s *Service) GetExtendedOrder(ctx context.Context, id int) (*models.ExtendedOrder, error) {
	if eo, ok := s.cache.Get(id); ok {
		s.log.Info("order loaded from cache", zap.Int("id", id))
		return eo, nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	extendedOrder := new(models.ExtendedOrder)

	order, err := s.ordersRepo.Get(ctx, tx, id)
	if err != nil {
		s.log.Error("failed to get order", zap.Error(err))
		return nil, err
	}
	extendedOrder.Order = *order

	items, err := s.itemsRepo.GetItems(ctx, tx, extendedOrder.Order.ID)
	if err != nil {
		s.log.Error("failed to get items", zap.Error(err))
		return nil, err
	}
	extendedOrder.Items = items

	payment, err := s.paymentRepo.Get(ctx, tx, extendedOrder.Order.PaymentID)
	if err != nil {
		s.log.Error("failed to get payment", zap.Error(err))
		return nil, err
	}
	extendedOrder.Payment = *payment

	delivery, err := s.deliveryRepo.Get(ctx, tx, extendedOrder.Order.DeliveryID)
	if err != nil {
		s.log.Error("failed to get delivery", zap.Error(err))
		return nil, err
	}
	extendedOrder.Delivery = *delivery

	if err = tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return nil, err
	}

	s.log.Info("order loaded from db", zap.Int("id", id))

	return extendedOrder, nil
}
