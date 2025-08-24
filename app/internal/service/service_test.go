package service

import (
	"test-task/internal/mocks"
	"test-task/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockExtendedOrderRepository(ctrl)

	service := NewService(nil, mockRepo, 10, zap.NewNop())

	id := 123
	eo := &models.ExtendedOrder{Order: models.Order{ID: id}}

	mockRepo.EXPECT().
		CreateExtendedOrder(gomock.Any(), eo).
		Return(nil)

	err := service.CreateExtendedOrder(t.Context(), eo)

	assert.NoError(t, err)
}

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockExtendedOrderRepository(ctrl)

	service := NewService(nil, mockRepo, 10, zap.NewNop())

	id := 123
	expected := &models.ExtendedOrder{Order: models.Order{ID: id}}

	mockRepo.EXPECT().
		GetExtendedOrder(gomock.Any(), id).
		Return(expected, nil)

	eo, err := service.GetExtendedOrder(t.Context(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected, eo)
}

func TestService_GetFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockExtendedOrderRepository(ctrl)

	service := NewService(nil, mockRepo, 10, zap.NewNop())

	id := 123
	expected := &models.ExtendedOrder{Order: models.Order{ID: id}}

	service.cache.Add(id, expected)

	eo, err := service.GetExtendedOrder(t.Context(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected, eo)
}

func TestService_LoadCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockExtendedOrderRepository(ctrl)

	maxCacheSize := 10
	limit := 5

	service := NewService(nil, mockRepo, maxCacheSize, zap.NewNop())

	expected := []*models.ExtendedOrder{
		{Order: models.Order{ID: 123}},
		{Order: models.Order{ID: 456}},
		{Order: models.Order{ID: 789}},
		{Order: models.Order{ID: 101112}},
		{Order: models.Order{ID: 131415}},
	}

	mockRepo.EXPECT().
		GetLastExtendedOrders(gomock.Any(), limit).
		Return(expected, nil)

	err := service.LoadRecentOrdersToCache(t.Context(), limit)

	assert.NoError(t, err)

	assert.Equal(t, len(expected), service.cache.Len())

	for _, eo := range expected {
		val, ok := service.cache.Get(eo.Order.ID)
		assert.True(t, ok, "Expected order to be in cache")
		assert.Equal(t, eo, val, "Expected order to match cache value")
	}
}
