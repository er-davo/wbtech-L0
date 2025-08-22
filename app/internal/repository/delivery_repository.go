package repository

import (
	"context"
	"test-task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository interface {
	Create(ctx context.Context, tx pgx.Tx, delivery *models.Delivery) error
	Get(ctx context.Context, tx pgx.Tx, id int) (*models.Delivery, error)
	GetByOrderIDs(ctx context.Context, tx pgx.Tx, ids []int) ([]*models.Delivery, error)
	Update(ctx context.Context, tx pgx.Tx, delivery *models.Delivery) error
	Delete(ctx context.Context, tx pgx.Tx, id int64) error
}

type deliveryRepository struct {
	db *pgxpool.Pool
}

func NewDeliveryRepository(db *pgxpool.Pool) DeliveryRepository {
	return &deliveryRepository{db: db}
}

func (r *deliveryRepository) Create(ctx context.Context, tx pgx.Tx, delivery *models.Delivery) error {
	if delivery == nil {
		return ErrNilValue
	}

	var exec pgx.Row
	if tx != nil {
		exec = tx.QueryRow(ctx, insertDeliveryQuery,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
		)
	} else {
		exec = r.db.QueryRow(ctx, insertDeliveryQuery,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
		)
	}

	err := exec.Scan(&delivery.ID)

	return wrapDBError(err)
}

func (r *deliveryRepository) Get(ctx context.Context, tx pgx.Tx, id int) (*models.Delivery, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	query := `
		SELECT
			id, name, phone, zip, city, address, region, email
		FROM delivery
		WHERE id = $1;
	`

	delivery := new(models.Delivery)
	var exec pgx.Row
	if tx != nil {
		exec = tx.QueryRow(ctx, query, id)
	} else {
		exec = r.db.QueryRow(ctx, query, id)
	}

	err := exec.Scan(
		&delivery.ID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	return delivery, wrapDBError(err)
}

func (r *deliveryRepository) GetByOrderIDs(ctx context.Context, tx pgx.Tx, ids []int) ([]*models.Delivery, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `
        SELECT
			id, name, phone, zip, city, address, region, email
        FROM delivery
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

	deliveries := make([]*models.Delivery, 0, len(ids))
	for rows.Next() {
		d := new(models.Delivery)
		if err := rows.Scan(
			&d.ID, &d.Name, &d.Phone, &d.Zip, &d.City,
			&d.Address, &d.Region, &d.Email,
		); err != nil {
			return nil, wrapDBError(err)
		}
		deliveries = append(deliveries, d)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return deliveries, nil
}

func (r *deliveryRepository) Update(ctx context.Context, tx pgx.Tx, delivery *models.Delivery) error {
	if delivery == nil || delivery.ID <= 0 {
		return ErrNilValue
	}

	query := `
		UPDATE delivery SET
			name = $1, phone = $2, zip = $3, city = $4, address = $5, region = $6, email = $7
		WHERE id = $8;
	`

	var cmd pgconn.CommandTag
	var err error
	if tx != nil {
		cmd, err = tx.Exec(ctx, query,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
			delivery.ID,
		)
	} else {
		cmd, err = r.db.Exec(ctx, query,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
			delivery.ID,
		)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return wrapDBError(err)
}

func (r *deliveryRepository) Delete(ctx context.Context, tx pgx.Tx, id int64) error {
	if id <= 0 {
		return ErrInvalidID
	}

	query := `DELETE FROM delivery WHERE id = $1;`

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
