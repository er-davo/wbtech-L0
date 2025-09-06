//go:build integration
// +build integration

package repository_test

import (
	"test-task/internal/models"
	"test-task/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtendedOrderRepository(t *testing.T) {
	repo := repository.NewExtendedOrderRepository(db)

	extendedOrder := &models.ExtendedOrder{
		Order: models.Order{
			OrderUID:        "extended order test",
			TrackNumber:     "2634",
			Entry:           "142",
			Locale:          "ru",
			CustomerID:      "test",
			DeliveryService: "test",
			ShardKey:        "test",
			SMID:            2,
			DateCreated:     time.Date(2025, time.September, 4, 3, 0, 0, 0, time.Local).UTC(),
			OOFShard:        "test",
		},
		Payment: models.Payment{
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
		},
		Delivery: models.Delivery{
			Name:    "test",
			Phone:   "+7926",
			Zip:     "1542",
			City:    "Moscow",
			Address: "Lenina",
			Region:  "Moscow",
			Email:   "test@emal.com",
		},
		Items: []*models.Item{
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
		},
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.CreateExtendedOrder(t.Context(), extendedOrder)
		assert.NoError(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		eo, err := repo.GetExtendedOrder(t.Context(), extendedOrder.Order.ID)
		assert.NoError(t, err)
		assert.Equal(t, extendedOrder, eo)
	})

	t.Run("Get Last N", func(t *testing.T) {
		eos, err := repo.GetLastExtendedOrders(t.Context(), 10)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(eos))
		assert.Equal(t, extendedOrder, eos[0])
	})
}
