//go:build integration
// +build integration

package repository_test

import (
	"testing"

	"test-task/internal/models"
	"test-task/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestPaymentRepository_CRUD(t *testing.T) {
	repo := repository.NewPaymentRepository(db)

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

	tx, err := db.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(t.Context())

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(t.Context(), tx, payment)
		assert.NoError(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		p, err := repo.Get(t.Context(), tx, payment.ID)
		assert.Equal(t, payment, p)
		assert.NoError(t, err)
	})

	payment.Transaction = "new test"

	t.Run("Update", func(t *testing.T) {
		err := repo.Update(t.Context(), tx, payment)
		assert.NoError(t, err)
		p, err := repo.Get(t.Context(), tx, payment.ID)
		assert.Equal(t, payment, p)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(t.Context(), tx, payment.ID)
		assert.NoError(t, err)
		_, err = repo.Get(t.Context(), tx, payment.ID)
		assert.ErrorIs(t, repository.ErrNotFound, err)
	})
}
