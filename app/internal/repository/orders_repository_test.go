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

func TestOrdersRepository_CRUD(t *testing.T) {
	dRepo := repository.NewDeliveryRepository(db)
	pRepo := repository.NewPaymentRepository(db)
	repo := repository.NewOrdersRepository(db)

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
		OrderUID:        "orders crud test",
		TrackNumber:     "2634",
		Entry:           "142",
		Locale:          "ru",
		CustomerID:      "test",
		DeliveryService: "test",
		ShardKey:        "test",
		SMID:            2,
		DateCreated:     time.Date(2025, time.September, 2, 0, 0, 0, 0, time.Local).UTC(),
		OOFShard:        "test",
	}

	tx, err := db.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(t.Context())

	t.Run("Create", func(t *testing.T) {
		err := dRepo.Create(t.Context(), tx, delivery)
		assert.NoError(t, err)
		err = pRepo.Create(t.Context(), tx, payment)
		assert.NoError(t, err)

		order.DeliveryID = delivery.ID
		order.PaymentID = payment.ID

		t.Logf("delivery id: %d, payment id: %d\n", delivery.ID, payment.ID)
		t.Logf("order: delivery id: %d, payment id: %d\n", order.DeliveryID, order.PaymentID)

		err = repo.Create(t.Context(), tx, order)
		assert.NoError(t, err)
		t.Logf("order id: %d\n", order.ID)
	})

	t.Run("Get", func(t *testing.T) {
		o, err := repo.Get(t.Context(), tx, order.ID)
		assert.Equal(t, order, o)
		assert.NoError(t, err)
	})

	order.DeliveryService = "new test"

	t.Run("Update", func(t *testing.T) {
		err := repo.Update(t.Context(), tx, order)
		assert.NoError(t, err)
		o, err := repo.Get(t.Context(), tx, order.ID)
		assert.Equal(t, order, o)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(t.Context(), tx, order.ID)
		assert.NoError(t, err)
		_, err = repo.Get(t.Context(), tx, order.ID)
		assert.ErrorIs(t, repository.ErrNotFound, err)
	})
}
