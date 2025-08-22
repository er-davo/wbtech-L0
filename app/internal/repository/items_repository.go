package repository

import (
	"context"

	"test-task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemsRepository interface {
	CreateItem(ctx context.Context, tx pgx.Tx, item *models.Item) error
	CreateItems(ctx context.Context, tx pgx.Tx, items []*models.Item) error
	Get(ctx context.Context, tx pgx.Tx, id int) (*models.Item, error)
	GetItems(ctx context.Context, tx pgx.Tx, orderID int) ([]*models.Item, error)
	GetByOrderIDs(ctx context.Context, tx pgx.Tx, orderIDs []int) ([]*models.Item, error)
	Update(ctx context.Context, tx pgx.Tx, item *models.Item) error
	Delete(ctx context.Context, tx pgx.Tx, id int) error
}

type itemsRepository struct {
	db *pgxpool.Pool
}

func NewItemsRepository(db *pgxpool.Pool) ItemsRepository {
	return &itemsRepository{
		db: db,
	}
}

func (r *itemsRepository) CreateItem(ctx context.Context, tx pgx.Tx, item *models.Item) error {
	if item == nil {
		return ErrNilValue
	}

	var err error
	if tx == nil {
		err = r.db.QueryRow(
			ctx,
			insertItemQuery,
			item.OrderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		).Scan(item.ID)
	} else {
		err = tx.QueryRow(
			ctx,
			insertItemQuery,
			item.OrderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		).Scan(item.ID)
	}

	return wrapDBError(err)
}

func (r *itemsRepository) CreateItems(ctx context.Context, tx pgx.Tx, items []*models.Item) error {
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		if item == nil {
			return ErrNilValue
		}
	}

	batch := &pgx.Batch{}
	for _, item := range items {
		batch.Queue(
			insertItemQuery,
			item.OrderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		)
	}

	if tx == nil {
		br := r.db.SendBatch(ctx, batch)
		defer br.Close()

		for _, item := range items {
			if err := br.QueryRow().Scan(&item.ID); err != nil {
				return wrapDBError(err)
			}
		}
	} else {
		br := tx.SendBatch(ctx, batch)
		defer br.Close()

		for _, item := range items {
			if err := br.QueryRow().Scan(&item.ID); err != nil {
				return wrapDBError(err)
			}
		}
	}

	return nil
}

func (r *itemsRepository) Get(ctx context.Context, tx pgx.Tx, id int) (*models.Item, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	item := new(models.Item)
	item.ID = id
	query := selectItemsWitoutWhereQuery + `
		WHERE id = $1;
	`

	var err error
	if tx == nil {
		err = r.db.QueryRow(
			ctx,
			query,
			item.ID,
		).Scan(
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
	} else {
		err = tx.QueryRow(
			ctx,
			query,
			item.ID,
		).Scan(
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
	}

	return item, wrapDBError(err)
}

func (r *itemsRepository) GetItems(ctx context.Context, tx pgx.Tx, orderID int) ([]*models.Item, error) {
	if orderID <= 0 {
		return nil, ErrInvalidID
	}

	query := selectItemsWitoutWhereQuery + `
		WHERE order_id = $1;
	`

	var rows pgx.Rows
	var err error
	if tx == nil {
		rows, err = r.db.Query(ctx, query, orderID)
	} else {
		rows, err = tx.Query(ctx, query, orderID)
	}
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	var items []*models.Item
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
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrNotFound
	}

	return items, nil
}

func (r *itemsRepository) GetByOrderIDs(ctx context.Context, tx pgx.Tx, orderIDs []int) ([]*models.Item, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}

	query := `
        SELECT
			id,
			order_id,
			chrt_id,
			track_number,
			price,
			rid,
			name,
			sale,
			size,
			total_price,
			nm_id,
			brand,
			status
		FROM items
        WHERE order_id = ANY($1);
    `

	var rows pgx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(ctx, query, orderIDs)
	} else {
		rows, err = r.db.Query(ctx, query, orderIDs)
	}
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	items := make([]*models.Item, 0, len(orderIDs))
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
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return items, nil
}

func (r *itemsRepository) Update(ctx context.Context, tx pgx.Tx, item *models.Item) error {
	if item == nil {
		return ErrNilValue
	}
	if item.ID <= 0 {
		return ErrInvalidID
	}

	query := `
		UPDATE items
		SET
			order_id = $2,
			chrt_id = $3,
			track_number = $4,
			price = $5,
			rid = $6,
			name = $7,
			sale = $8,
			size = $9,
			total_price = $10,
			nm_id = $11,
			brand = $12,
			status = $13
		WHERE id = $1;
	`

	var cmd pgconn.CommandTag
	var err error
	if tx == nil {
		cmd, err = r.db.Exec(
			ctx,
			query,
			item.ID,
			item.OrderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		)
	} else {
		cmd, err = tx.Exec(
			ctx,
			query,
			item.ID,
			item.OrderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return wrapDBError(err)
}

func (r *itemsRepository) Delete(ctx context.Context, tx pgx.Tx, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	query := `
		DELETE FROM items
		WHERE id = $1;
	`

	var cmd pgconn.CommandTag
	var err error
	if tx == nil {
		cmd, err = r.db.Exec(ctx, query, id)
	} else {
		cmd, err = tx.Exec(ctx, query, id)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return wrapDBError(err)
}
