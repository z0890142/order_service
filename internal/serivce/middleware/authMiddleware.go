package middleware

import (
	"net/http"
	"order_service/config"
	"order_service/internal/data"
	"order_service/pkg/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type middleware struct {
	dataMgr data.DataManager
}

func NewMiddleware(dataMgr data.DataManager) *middleware {
	return &middleware{
		dataMgr: dataMgr,
	}
}

func (m *middleware) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.Split(authHeader, " ")[1]

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Jwt.Key), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	doctor := models.Doctor{}
	// Check if Doctor exists with the given ID (you need to implement this logic)
	err = m.dataMgr.GetDoctor(c, map[string]interface{}{
		"id": claims.Id,
	}, &doctor)
	if err != nil || doctor.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Set Doctor in context for further use in handlers
	c.Set("doctor", doctor)

	c.Next()
}
