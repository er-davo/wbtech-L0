//go:build integration
// +build integration

package repository_test

import (
	"testing"

	"test-task/internal/models"
	"test-task/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestDeliveryRepository_CRUD(t *testing.T) {
	repo := repository.NewDeliveryRepository(db)

	delivery := &models.Delivery{
		Name:    "test",
		Phone:   "+7926",
		Zip:     "1542",
		City:    "Moscow",
		Address: "Lenina",
		Region:  "Moscow",
		Email:   "test@emal.com",
	}

	tx, err := db.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(t.Context())

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(t.Context(), tx, delivery)
		assert.NotEqual(t, 0, delivery.ID)
		assert.NoError(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		d, err := repo.Get(t.Context(), tx, delivery.ID)
		assert.Equal(t, delivery, d)
		assert.NoError(t, err)
	})

	delivery.Name = "new test"

	t.Run("Update", func(t *testing.T) {
		err := repo.Update(t.Context(), tx, delivery)
		assert.NoError(t, err)
		d, err := repo.Get(t.Context(), tx, delivery.ID)
		assert.Equal(t, delivery, d)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(t.Context(), tx, delivery.ID)
		assert.NoError(t, err)
		_, err = repo.Get(t.Context(), tx, delivery.ID)
		assert.ErrorIs(t, repository.ErrNotFound, err)
	})
}
