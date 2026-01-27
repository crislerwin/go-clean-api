package handler

import (
	"log/slog"
	"net/http"

	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	loginUseCase *usecase.LoginUseCase
}

func NewAuthHandler(loginUseCase *usecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
	}
}

// Login godoc
// @Summary      Authenticate user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body LoginRequest true "User Credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid login request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	input := usecase.LoginInputDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		slog.Warn("Login failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, output)
}
