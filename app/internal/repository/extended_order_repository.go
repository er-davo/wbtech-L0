package repository

import (
	"context"

	"test-task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExtendedOrderRepository interface {
	CreateExtendedOrder(ctx context.Context, eo *models.ExtendedOrder) error
	GetExtendedOrder(ctx context.Context, id int) (*models.ExtendedOrder, error)
	Orders() OrdersRepository
	Items() ItemsRepository
	Delivery() DeliveryRepository
	Payment() PaymentRepository
}

type extendedOrderRepository struct {
	db       *pgxpool.Pool
	orders   OrdersRepository
	items    ItemsRepository
	delivery DeliveryRepository
	payment  PaymentRepository
}

func NewExtendedOrderRepository(db *pgxpool.Pool) ExtendedOrderRepository {
	return &extendedOrderRepository{
		db:       db,
		orders:   NewOrdersRepository(db),
		items:    NewItemsRepository(db),
		delivery: NewDeliveryRepository(db),
		payment:  NewPaymentRepository(db),
	}
}

func (r *extendedOrderRepository) CreateExtendedOrder(ctx context.Context, eo *models.ExtendedOrder) error {
	if eo == nil {
		return ErrNilValue
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return wrapDBError(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = r.delivery.Create(ctx, tx, &eo.Delivery)
	if err != nil {
		return wrapDBError(err)
	}

	err = r.payment.Create(ctx, tx, &eo.Payment)
	if err != nil {
		return wrapDBError(err)
	}

	eo.Order.DeliveryID = eo.Delivery.ID
	eo.Order.PaymentID = eo.Payment.ID

	err = r.orders.Create(ctx, tx, &eo.Order)
	if err != nil {
		return wrapDBError(err)
	}

	for _, item := range eo.Items {
		item.OrderID = eo.Order.ID
	}

	err = r.items.CreateItems(ctx, tx, eo.Items)
	if err != nil {
		return wrapDBError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapDBError(err)
	}

	return nil
}

func (r *extendedOrderRepository) GetExtendedOrder(ctx context.Context, id int) (*models.ExtendedOrder, error) {
	if id < 0 {
		return nil, ErrInvalidID
	}

	eo := new(models.ExtendedOrder)

	batch := &pgx.Batch{}

	query := selectExtendedOrderWithoutItemsQuery + `
		WHERE o.id = $1;
	`
	itemsQuery := selectItemsWitoutWhereQuery + `
		wHERE order_id = $1;
	`

	batch.Queue(
		query,
		id,
	)

	batch.Queue(
		itemsQuery,
		id,
	)

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	err := br.QueryRow().Scan(
		&eo.Order.ID, &eo.Order.OrderUID, &eo.Order.TrackNumber,
		&eo.Order.Entry, &eo.Order.DeliveryID, &eo.Order.PaymentID,
		&eo.Order.Locale, &eo.Order.InternalSignature,
		&eo.Order.CustomerID, &eo.Order.DeliveryService,
		&eo.Order.ShardKey, &eo.Order.SMID, &eo.Order.DateCreated, &eo.Order.OOFShard,

		&eo.Delivery.ID, &eo.Delivery.Name, &eo.Delivery.Phone, &eo.Delivery.Zip, &eo.Delivery.City, &eo.Delivery.Address,

		&eo.Payment.ID, &eo.Payment.Transaction, &eo.Payment.RequestID,
		&eo.Payment.Currency, &eo.Payment.Provider, &eo.Payment.Amount,
		&eo.Payment.PaymentDate, &eo.Payment.Bank, &eo.Payment.DeliveryCost,
		&eo.Payment.GoodsTotal, &eo.Payment.CustomFee,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}

	rows, err := br.Query()
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	eo.Items = make([]*models.Item, 0)
	for rows.Next() {
		item := new(models.Item)
		err = rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, wrapDBError(err)
		}
		eo.Items = append(eo.Items, item)
	}

	return eo, nil
}

func (r *extendedOrderRepository) GetLastExtendedOrders(ctx context.Context, limit int) ([]*models.ExtendedOrder, error) {
	if limit < 0 {
		return nil, ErrNilValue
	}

	if limit == 0 {
		return []*models.ExtendedOrder{}, nil
	}

	query := selectExtendedOrderWithoutItemsQuery + `
		ORDER BY o.date_created DESC
		LIMIT $1;
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	eos := make([]*models.ExtendedOrder, 0, limit)

	for rows.Next() {
		eo := new(models.ExtendedOrder)
		if err := rows.Scan(
			&eo.Order.ID, &eo.Order.OrderUID, &eo.Order.TrackNumber,
			&eo.Order.Entry, &eo.Order.DeliveryID, &eo.Order.PaymentID,
			&eo.Order.Locale, &eo.Order.InternalSignature,
			&eo.Order.CustomerID, &eo.Order.DeliveryService,
			&eo.Order.ShardKey, &eo.Order.SMID, &eo.Order.DateCreated, &eo.Order.OOFShard,
			&eo.Delivery.ID, &eo.Delivery.Name, &eo.Delivery.Phone, &eo.Delivery.Zip, &eo.Delivery.City, &eo.Delivery.Address,
			&eo.Payment.ID, &eo.Payment.Transaction, &eo.Payment.RequestID,
			&eo.Payment.Currency, &eo.Payment.Provider, &eo.Payment.Amount,
			&eo.Payment.PaymentDate, &eo.Payment.Bank, &eo.Payment.DeliveryCost,
			&eo.Payment.GoodsTotal, &eo.Payment.CustomFee,
		); err != nil {
			return nil, wrapDBError(err)
		}
		eos = append(eos, eo)
	}

	orderIDs := make([]int, 0, len(eos))
	for _, eo := range eos {
		orderIDs = append(orderIDs, eo.Order.ID)
	}

	itemsQuery := selectItemsWitoutWhereQuery + `
		WHERE order_id = ANY($1);
	`

	rows, err = r.db.Query(ctx, itemsQuery, orderIDs)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	items := make(map[int][]*models.Item)
	for rows.Next() {
		item := new(models.Item)
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		); err != nil {
			return nil, wrapDBError(err)
		}
		items[item.OrderID] = append(items[item.OrderID], item)
	}

	for _, eo := range eos {
		if its, ok := items[eo.Order.ID]; ok {
			eo.Items = its
		}
	}

	return eos, nil
}

func (r *extendedOrderRepository) Orders() OrdersRepository { return r.orders }

func (r *extendedOrderRepository) Items() ItemsRepository { return r.items }

func (r *extendedOrderRepository) Delivery() DeliveryRepository { return r.delivery }

func (r *extendedOrderRepository) Payment() PaymentRepository { return r.payment }
