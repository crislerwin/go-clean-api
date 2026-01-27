package handler

import (
	"log/slog"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Name         string  `json:"name" binding:"required"`
	Location     string  `json:"location" binding:"required"`
	Organization string  `json:"organization" binding:"required"`
	Rating       string  `json:"rating" binding:"required"`
	Date         string  `json:"date" binding:"required" example:"2025-10-10T00:00:00Z"`
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

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Create a new event (Admin only)
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body CreateEventRequest true "Event Data"
// @Success      201  {object}  usecase.CreateEventOutputDTO
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err)
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
		slog.Error("Error creating event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	slog.Info("Event created successfully", "event_id", output.ID, "event_name", output.Name)
	c.JSON(http.StatusCreated, output)
}
