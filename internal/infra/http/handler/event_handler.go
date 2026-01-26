package handler

import (
	"log"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type createEventRequest struct {
	Name         string  `json:"name" binding:"required"`
	Location     string  `json:"location" binding:"required"`
	Organization string  `json:"organization" binding:"required"`
	Rating       string  `json:"rating" binding:"required"`
	Date         string  `json:"date" binding:"required"`
	Capacity     int     `json:"capacity" binding:"required,min=1"`
	Price        float64 `json:"price" binding:"required,min=0"`
	ImageURL     string  `json:"image_url" binding:"required"`
}

type EventHandler struct {
	createEventUseCase *usecase.CreateEventUseCase
}

func NewEventHandler(createEventUseCase *usecase.CreateEventUseCase) *EventHandler {
	return &EventHandler{
		createEventUseCase: createEventUseCase,
	}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req createEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	input := usecase.CreateEventInputDTO{
		Name:         req.Name,
		Location:     req.Location,
		Organization: req.Organization,
		Rating:       req.Rating,
		Date:         req.Date,
		Capacity:     req.Capacity,
		Price:        req.Price,
		ImageURL:     req.ImageURL,
	}

	output, err := h.createEventUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		log.Println("Error creating event:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, output)
}
