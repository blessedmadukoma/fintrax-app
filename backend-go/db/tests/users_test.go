package db_test

import (
	"context"
	db "fintrax/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	arg := db.CreateUserParams{
		Email:          "test@example4.com",
		HashedPassword: "secret",
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)

	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
}
