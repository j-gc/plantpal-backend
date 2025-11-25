package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j-gc/plantpal-backend/internal/modules/auth/application"
)

type AuthHandlers struct {
	svc *application.Service
}

func NewAuthHandlers(svc *application.Service) *AuthHandlers {
	return &AuthHandlers{svc: svc}
}

func (h *AuthHandlers) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.register)
	r.POST("/login", h.login)
}

func (h *AuthHandlers) register(c *gin.Context) {
	var in application.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	out, err := h.svc.Register(c, in)
	if err != nil {
		status := http.StatusBadRequest
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *AuthHandlers) login(c *gin.Context) {
	var in application.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	out, err := h.svc.Login(c, in)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, out)
}
