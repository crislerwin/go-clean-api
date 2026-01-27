package handler

import (
	"log/slog"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
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
	createEventUseCase    *usecase.CreateEventUseCase
	listEventsUseCase     *usecase.ListEventsUseCase
	listUserEventsUseCase *usecase.ListUserEventsUseCase
}

func NewEventHandler(
	createEventUseCase *usecase.CreateEventUseCase,
	listEventsUseCase *usecase.ListEventsUseCase,
	listUserEventsUseCase *usecase.ListUserEventsUseCase,
) *EventHandler {
	return &EventHandler{
		createEventUseCase:    createEventUseCase,
		listEventsUseCase:     listEventsUseCase,
		listUserEventsUseCase: listUserEventsUseCase,
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

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	output, err := h.createEventUseCase.Execute(c.Request.Context(), userID.(string), input)
	if err != nil {
		if err == entity.ErrDateInPast || err == entity.ErrInvalidEventData {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Error("Error creating event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	slog.Info("Event created successfully", "event_id", output.ID, "event_name", output.Name)
	c.JSON(http.StatusCreated, output)
}

// ListEvents godoc
// @Summary      List all events
// @Description  List all available events
// @Tags         events
// @Accept       json
// @Produce      json
// @Success      200  {array}   usecase.ListEventsOutputDTO
// @Failure      500  {object}  map[string]string
// @Router       /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	output, err := h.listEventsUseCase.Execute(c.Request.Context())
	if err != nil {
		slog.Error("Error listing events", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// ListMyEvents godoc
// @Summary      List my events
// @Description  List events created by the logged-in user
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   usecase.ListEventsOutputDTO
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /me/events [get]
func (h *EventHandler) ListMyEvents(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	output, err := h.listUserEventsUseCase.Execute(c.Request.Context(), userID.(string))
	if err != nil {
		slog.Error("Error listing user events", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, output)
}
