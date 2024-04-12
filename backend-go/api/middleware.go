package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			ctx.Abort()
			return
		}

		splitToken := strings.Split(token, " ")

		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			ctx.Abort()
			return
		}

		userId, err := tokenController.VerifyToken(splitToken[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userId)
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			// log.Println("got options and stopped")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
