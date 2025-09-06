package repository

import (
	"context"
	"test-task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersRepository interface {
	Create(ctx context.Context, tx pgx.Tx, order *models.Order) error
	Get(ctx context.Context, tx pgx.Tx, id int64) (*models.Order, error)
	Update(ctx context.Context, tx pgx.Tx, order *models.Order) error
	Delete(ctx context.Context, tx pgx.Tx, id int64) error
}

type ordersRepository struct {
	db *pgxpool.Pool
}

func NewOrdersRepository(db *pgxpool.Pool) OrdersRepository {
	return &ordersRepository{db: db}
}

func (r *ordersRepository) Create(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	if order == nil {
		return ErrNilValue
	}

	var exec pgx.Row
	if tx != nil {
		exec = tx.QueryRow(ctx, insertOrderQuery,
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.DeliveryID,
			order.PaymentID,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SMID,
			order.DateCreated,
			order.OOFShard,
		)
	} else {
		exec = r.db.QueryRow(ctx, insertOrderQuery,
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.DeliveryID,
			order.PaymentID,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SMID,
			order.DateCreated,
			order.OOFShard,
		)
	}

	err := exec.Scan(&order.ID)

	return wrapDBError(err)
}

func (r *ordersRepository) Get(ctx context.Context, tx pgx.Tx, id int64) (*models.Order, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	query := `
		SELECT
			id,
			order_uid,
			track_number,
			entry,
			delivery_id,
			payment_id,
			locale,
			internal_signature,
			customer_id,
			delivery_service,
			shardkey,
			sm_id,
			date_created,
			oof_shard
		FROM orders
		WHERE id = $1;
	`

	order := new(models.Order)
	var exec pgx.Row
	if tx != nil {
		exec = tx.QueryRow(ctx, query, id)
	} else {
		exec = r.db.QueryRow(ctx, query, id)
	}

	err := exec.Scan(
		&order.ID,
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.DeliveryID,
		&order.PaymentID,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SMID,
		&order.DateCreated,
		&order.OOFShard,
	)

	return order, wrapDBError(err)
}

func (r *ordersRepository) Update(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	if order == nil {
		return ErrNilValue
	}

	query := `
		UPDATE orders SET
			order_uid = $2,
			track_number = $3,
			entry = $4,
			delivery_id = $5,
			payment_id = $6,
			locale = $7,
			internal_signature = $8,
			customer_id = $9,
			delivery_service = $10,
			shardkey = $11,
			sm_id = $12,
			oof_shard = $13
		WHERE id = $1;
	`

	var cmd pgconn.CommandTag
	var err error
	if tx != nil {
		cmd, err = tx.Exec(ctx, query,
			order.ID,
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.DeliveryID,
			order.PaymentID,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SMID,
			order.OOFShard,
		)
	} else {
		cmd, err = r.db.Exec(ctx, query,
			order.ID,
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.DeliveryID,
			order.PaymentID,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SMID,
			order.OOFShard,
		)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return wrapDBError(err)
}

func (r *ordersRepository) Delete(ctx context.Context, tx pgx.Tx, id int64) error {
	if id <= 0 {
		return ErrInvalidID
	}

	query := `DELETE FROM orders WHERE id = $1;`
	var cmd pgconn.CommandTag
	var err error
	if tx != nil {
		cmd, err = tx.Exec(ctx, query, id)
	} else {
		cmd, err = r.db.Exec(ctx, query, id)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return wrapDBError(err)
}
