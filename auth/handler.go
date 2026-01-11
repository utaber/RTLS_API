package auth

import (
	"RTLS_API/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	jwtService  *JWTService
	userService UserService
}

func NewHandler(jwt *JWTService, userSvc UserService) *Handler {
	return &Handler{
		jwtService:  jwt,
		userService: userSvc,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	username, err := h.userService.AuthenticateByEmail(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := h.jwtService.GenerateToken(username)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token": token,
	})
}
