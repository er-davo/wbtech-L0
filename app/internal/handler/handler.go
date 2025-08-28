package handler

import (
	"errors"
	"net/http"
	"strconv"
	"test-task/internal/models"
	"test-task/internal/repository"
	"test-task/internal/retry"
	"test-task/internal/service"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Handler struct {
	service *service.Service
	retry   retry.Retrier
	log     *zap.Logger
}

func NewHandler(service *service.Service, retry retry.Retrier, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		retry:   retry,
		log:     log,
	}
}

func (h *Handler) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID format",
		})
	}

	eo := new(models.ExtendedOrder)

	h.log.Info("geting order", zap.Int("id", id))

	if err := h.retry.Do(c.Request().Context(), func(attempt int) error {
		var err error
		if eo, err = h.service.GetExtendedOrder(c.Request().Context(), id); err != nil {
			h.log.Warn("error on getting order", zap.Int("id", id), zap.Error(err), zap.Int("attempt", attempt))
			return err
		}
		h.log.Info("order found", zap.Int("id", id))
		return nil
	}); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.log.Warn("order not found", zap.Int("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Order not found"})
		} else {
			h.log.Error("error on getting order", zap.Int("id", id), zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}
	}

	return c.JSON(http.StatusOK, eo)
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/order")
	g.GET("/:id", h.Get)
}
