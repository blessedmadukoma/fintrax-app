package api

import (
	"database/sql"
	"fintrax/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	server *Server
}

func (a Auth) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/auth")

	serverGroup.POST("login", a.login)
}

func (a Auth) login(ctx *gin.Context) {
	user := new(UserParams)

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error binding data": err.Error()})
		return
	}

	dbUser, err := a.server.queries.GetUserByEmail(ctx, user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("user not found: %v", err.Error())})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.VerifyPassword(user.Password, dbUser.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid password: %v", err.Error())})
		return
	}

	token, err := utils.CreateToken(dbUser.ID, a.server.config.SIGNINGKEY)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create token: %v", err.Error())})
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
