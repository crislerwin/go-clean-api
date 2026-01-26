package handler

import (
	"fmt"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	createOrderUseCase *usecase.CreateOrderUseCase
}

func NewOrderHandler(createOrderUseCase *usecase.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input usecase.OrderInputDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	output, err := h.createOrderUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch err {
		case usecase.ErrEventSoldOut:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case usecase.ErrEventNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			fmt.Println("Error creating order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, output)

}
