package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (j *JWTService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Split(auth, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"detail": "Invalid auth"})
			return
		}

		_, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(j.Secret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"detail": "Invalid token"})
			return
		}

		c.Next()
	}
}
