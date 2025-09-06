//go:build integration
// +build integration

package repository_test

import (
	"testing"
	"time"

	"test-task/internal/models"
	"test-task/internal/repository"

	"github.com/stretchr/testify/assert"
)

// test build and change int64 to int, delete last migration

func TestItemsRepository_CRUD(t *testing.T) {
	dRepo := repository.NewDeliveryRepository(db)
	pRepo := repository.NewPaymentRepository(db)
	oRepo := repository.NewOrdersRepository(db)
	repo := repository.NewItemsRepository(db)

	payment := &models.Payment{
		Transaction:  "test",
		RequestID:    "",
		Currency:     "RUB",
		Provider:     "alfa",
		Amount:       1000,
		PaymentDate:  90872534,
		Bank:         "tbank",
		DeliveryCost: 325,
		GoodsTotal:   32,
		CustomFee:    0,
	}

	delivery := &models.Delivery{
		Name:    "test",
		Phone:   "+7926",
		Zip:     "1542",
		City:    "Moscow",
		Address: "Lenina",
		Region:  "Moscow",
		Email:   "test@emal.com",
	}

	order := &models.Order{
		OrderUID:        "items crud test",
		TrackNumber:     "2634",
		Entry:           "142",
		Locale:          "ru",
		CustomerID:      "test",
		DeliveryService: "test",
		ShardKey:        "test",
		SMID:            2,
		DateCreated:     time.Date(2025, time.September, 4, 3, 0, 0, 0, time.Local).UTC(),
		OOFShard:        "test",
	}

	item := &models.Item{
		ChrtID:      324,
		TrackNumber: "test",
		Price:       200,
		RID:         "test",
		Name:        "test",
		Sale:        20,
		Size:        "test",
		TotalPrice:  200,
		NMID:        12,
		Brand:       "test",
		Status:      1,
	}

	items := []*models.Item{
		{
			ChrtID:      324,
			TrackNumber: "test",
			Price:       200,
			RID:         "test",
			Name:        "test",
			Sale:        20,
			Size:        "test",
			TotalPrice:  200,
			NMID:        12,
			Brand:       "test",
			Status:      1,
		},
		{
			ChrtID:      324,
			TrackNumber: "test",
			Price:       200,
			RID:         "test",
			Name:        "test",
			Sale:        20,
			Size:        "test",
			TotalPrice:  200,
			NMID:        12,
			Brand:       "test",
			Status:      1,
		},
		{
			ChrtID:      324,
			TrackNumber: "test",
			Price:       200,
			RID:         "test",
			Name:        "test",
			Sale:        20,
			Size:        "test",
			TotalPrice:  200,
			NMID:        12,
			Brand:       "test",
			Status:      1,
		},
	}

	tx, err := db.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(t.Context())

	t.Run("Create item", func(t *testing.T) {
		err := dRepo.Create(t.Context(), tx, delivery)
		assert.NoError(t, err)
		err = pRepo.Create(t.Context(), tx, payment)
		assert.NoError(t, err)

		order.DeliveryID = delivery.ID
		order.PaymentID = payment.ID

		err = oRepo.Create(t.Context(), tx, order)
		assert.NoError(t, err)

		item.OrderID = order.ID

		err = repo.CreateItem(t.Context(), tx, item)
		assert.NoError(t, err)
	})

	t.Run("Get item", func(t *testing.T) {
		i, err := repo.Get(t.Context(), tx, item.ID)
		assert.Equal(t, item, i)
		assert.NoError(t, err)
	})

	item.Name = "new test"

	t.Run("Update", func(t *testing.T) {
		err := repo.Update(t.Context(), tx, item)
		assert.NoError(t, err)
		i, err := repo.Get(t.Context(), tx, item.ID)
		assert.Equal(t, item, i)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(t.Context(), tx, item.ID)
		assert.NoError(t, err)
		_, err = repo.Get(t.Context(), tx, item.ID)
		assert.ErrorIs(t, repository.ErrNotFound, err)
	})

	items[0].OrderID = order.ID
	items[1].OrderID = order.ID
	items[2].OrderID = order.ID

	t.Run("Create items", func(t *testing.T) {
		err := repo.CreateItems(t.Context(), tx, items)
		assert.NoError(t, err)
	})

	t.Run("Get items", func(t *testing.T) {
		i, err := repo.GetItems(t.Context(), tx, order.ID)
		assert.Equal(t, items, i)
		assert.NoError(t, err)
	})
}
