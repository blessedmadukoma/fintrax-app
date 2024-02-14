package db_test

import (
	"context"
	db "fintrax/db/sqlc"
	"fintrax/utils"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	hashedPassword, err := utils.GenerateHashPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate password:", err)
	}

	arg := db.CreateUserParams{
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)

	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)

	// create a new user using the same details and verify if the email is the same: it should return an error
	user2, err := testQuery.CreateUser(context.Background(), arg)
	assert.Error(t, err)
	assert.Empty(t, user2)
}
