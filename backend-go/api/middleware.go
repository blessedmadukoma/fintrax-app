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
