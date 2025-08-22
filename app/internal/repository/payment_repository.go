package repository

import (
	"context"
	"test-task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(ctx context.Context, tx pgx.Tx, payment *models.Payment) error
	Get(ctx context.Context, tx pgx.Tx, id int) (*models.Payment, error)
	GetByOrderIDs(ctx context.Context, tx pgx.Tx, ids []int) ([]*models.Payment, error)
	Update(ctx context.Context, tx pgx.Tx, payment *models.Payment) error
	Delete(ctx context.Context, tx pgx.Tx, id int) error
}

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, tx pgx.Tx, payment *models.Payment) error {
	if payment == nil {
		return ErrNilValue
	}

	var exec pgx.Row
	if tx != nil {
		exec = tx.QueryRow(ctx, insertPaymentQuery,
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDate,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
		)
	} else {
		exec = r.db.QueryRow(ctx, insertPaymentQuery,
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDate,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
		)
	}

	err := exec.Scan(&payment.ID)

	return wrapDBError(err)
}

func (r *paymentRepository) Get(ctx context.Context, tx pgx.Tx, id int) (*models.Payment, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	query := `
		SELECT
			id,
			transaction,
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee
		FROM payments
		WHERE id=$1;
	`

	payment := new(models.Payment)
	var exec pgx.Row
	if tx != nil {
		exec = tx.QueryRow(ctx, query, id)
	} else {
		exec = r.db.QueryRow(ctx, query, id)
	}

	err := exec.Scan(
		&payment.ID,
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	return payment, wrapDBError(err)
}

func (r *paymentRepository) GetByOrderIDs(ctx context.Context, tx pgx.Tx, ids []int) ([]*models.Payment, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `
        SELECT id, transaction, request_id, currency, provider, amount,
               payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payment
        WHERE id = ANY($1);
    `

	var rows pgx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(ctx, query, ids)
	} else {
		rows, err = r.db.Query(ctx, query, ids)
	}
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	payments := make([]*models.Payment, 0, len(ids))
	for rows.Next() {
		p := new(models.Payment)
		if err := rows.Scan(
			&p.ID, &p.Transaction, &p.RequestID,
			&p.Currency, &p.Provider, &p.Amount,
			&p.PaymentDate, &p.Bank, &p.DeliveryCost,
			&p.GoodsTotal, &p.CustomFee,
		); err != nil {
			return nil, wrapDBError(err)
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return payments, nil
}

func (r *paymentRepository) Update(ctx context.Context, tx pgx.Tx, payment *models.Payment) error {
	if payment == nil {
		return ErrNilValue
	}

	query := `
		UPDATE payments SET
			transaction = $1,
			request_id = $2,
			currency = $3,
			provider = $4,
			amount = $5,
			payment_dt = $6,
			bank = $7,
			delivery_cost = $8,
			goods_total = $9,
			custom_fee = $10
		WHERE id = $11;
	`

	var cmd pgconn.CommandTag
	var err error
	if tx != nil {
		cmd, err = tx.Exec(ctx, query,
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDate,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
			payment.ID,
		)
	} else {
		cmd, err = r.db.Exec(ctx, query,
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDate,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
			payment.ID,
		)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return wrapDBError(err)
}

func (r *paymentRepository) Delete(ctx context.Context, tx pgx.Tx, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	query := `DELETE FROM payments WHERE id = $1;`

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
