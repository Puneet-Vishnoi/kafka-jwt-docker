package middleware

import (
	"net/http"

	"github.com/Puneet-Vishnoi/kafka-simple/helpers"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader)<7 || authHeader[:7] != "Bearer "{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer {token}"})
			return 
		}

		tokenStr := authHeader[7:]
		claims, err := helpers.ValidateToken(tokenStr)
		if err != nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return 
		}

		c.Set("user_emaill", claims.Email)
		c.Set("user_uid", claims.Uid)
		c.Next()


	}
}