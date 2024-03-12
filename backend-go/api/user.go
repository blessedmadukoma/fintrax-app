package api

import (
	"context"
	db "fintrax/db/sqlc"
	"fintrax/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type User struct {
	server *Server
}

func (u User) router(server *Server) {
	u.server = server

	serverGroup := server.router.Group("/users")

	serverGroup.GET("", u.listUsers)
	serverGroup.POST("", u.createUser)
}

type UserParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toNewUserResponse(user *db.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *User) createUser(c *gin.Context) {
	var user UserParams

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}

	newUser, err := u.server.queries.CreateUser(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// violated unique constraint i.e. user already exists
			if pgErr.Code == "23505" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{}.toNewUserResponse(&newUser))
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
