package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/http/auth"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	EventID  string `json:"event_id" binding:"required,uuid"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type OrderHandler struct {
	createOrderUseCase       *usecase.CreateOrderUseCase
	listUserOrdersUseCase    *usecase.ListUserOrdersUseCase
	updateOrderStatusUseCase *usecase.UpdateOrderStatusUseCase
}

func NewOrderHandler(create *usecase.CreateOrderUseCase, list *usecase.ListUserOrdersUseCase, update *usecase.UpdateOrderStatusUseCase) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase:       create,
		listUserOrdersUseCase:    list,
		updateOrderStatusUseCase: update,
	}
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PAID REJECTED"`
}

// UpdateStatus godoc
// @Summary      Update order status (Webhook)
// @Description  Update order status (e.g. from payment gateway)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id   path      string                    true "Order ID"
// @Param        input body      UpdateOrderStatusRequest  true "Status Data"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders/{id}/status [post]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	// Simple Auth
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		slog.Warn("WEBHOOK_SECRET not set, allowing request (DANGEROUS)")
	} else {
		token := c.GetHeader("X-Webhook-Secret")
		if token != secret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	orderID := c.Param("id")
	var req UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateOrderStatusInputDTO{
		OrderID: orderID,
		Status:  req.Status,
	}

	err := h.updateOrderStatusUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		slog.Error("Error updating order status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateOrder godoc
// @Summary      Purchase tickets
// @Description  Purchase tickets for an event
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body CreateOrderRequest true "Order Data"
// @Success      201  {object}  usecase.OrderOutputDTO
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid order request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	input := usecase.OrderInputDTO{
		EventID:  req.EventID,
		UserID:   userID,
		Quantity: req.Quantity,
	}

	output, err := h.createOrderUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch err {
		case entity.ErrEventSoldOut:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case usecase.ErrEventNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			slog.Error("Error creating order", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	slog.Info("Order created successfully", "order_id", output.ID, "user_id", userID, "event_id", input.EventID)
	c.JSON(http.StatusCreated, output)
}

// ListMyOrders godoc
// @Summary      List my orders
// @Description  List all orders for the authenticated user
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   usecase.ListOrdersOutputDTO
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [get]
func (h *OrderHandler) ListMyOrders(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	orders, err := h.listUserOrdersUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		slog.Error("Error listing orders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
