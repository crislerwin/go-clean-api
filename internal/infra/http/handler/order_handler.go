package handler

import (
	"log/slog"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/infra/http/auth"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	EventID  string `json:"event_id" binding:"required,uuid"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type OrderHandler struct {
	createOrderUseCase *usecase.CreateOrderUseCase
}

func NewOrderHandler(createOrderUseCase *usecase.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
	}
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
		case usecase.ErrEventSoldOut:
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
