package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token, utils.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("userId", claims.Sub)
		c.Set("useEmail", claims.Email)

		c.Next()
	}
}
