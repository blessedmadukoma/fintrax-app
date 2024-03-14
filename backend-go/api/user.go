package api

import (
	"context"
	"database/sql"
	db "fintrax/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	server *Server
}

func (u User) router(server *Server) {
	u.server = server

	serverGroup := server.router.Group("/users", AuthenticatedMiddleware())

	serverGroup.GET("", u.listUsers)

	serverGroup.GET("me", u.gettLoggedInUser)
}

func (u *User) listUsers(c *gin.Context) {
	arg := db.ListUsersParams{
		Offset: 0,
		Limit:  10,
	}
	users, err := u.server.queries.ListUsers(context.Background(), arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUsers := []UserResponse{}

	for _, user := range users {
		n := UserResponse{}.toNewUserResponse(&user)
		newUsers = append(newUsers, *n)
	}

	c.JSON(http.StatusOK, newUsers)
}

func (u *User) gettLoggedInUser(ctx *gin.Context) {
	value, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to access resources"})
		return
	}

	userId, ok := value.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Encounted an issue"})
	}

	user, err := u.server.queries.GetUserByID(ctx, userId)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, UserResponse{}.toNewUserResponse(&user))
}
