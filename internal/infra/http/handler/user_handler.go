package handler

import (
	"log/slog"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserHandler struct {
	signUpUseCase *usecase.SignUpUseCase
}

func NewUserHandler(signUpUseCase *usecase.SignUpUseCase) *UserHandler {
	return &UserHandler{
		signUpUseCase: signUpUseCase,
	}
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Register a new user with the default 'user' role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body SignUpRequest true "User Registration Credentials"
// @Success      201  {object}  usecase.SignUpOutputDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid user request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	input := usecase.SignUpInputDTO{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.signUpUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		slog.Error("Error creating user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	slog.Info("User created successfully", "user_id", output.ID)
	c.JSON(http.StatusCreated, output)
}
