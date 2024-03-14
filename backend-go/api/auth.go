package api

import (
	"context"
	"database/sql"
	db "fintrax/db/sqlc"
	"fintrax/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Auth struct {
	server *Server
}

func (a Auth) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/auth")

	serverGroup.POST("login", a.login)
	serverGroup.POST("register", a.register)
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

	token, err := tokenController.CreateToken(dbUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create token: %v", err.Error())})
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
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

func (a *Auth) register(c *gin.Context) {
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

	newUser, err := a.server.queries.CreateUser(context.Background(), arg)
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
