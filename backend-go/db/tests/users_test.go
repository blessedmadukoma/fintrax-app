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

func createRandomUser(t *testing.T) db.User {
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
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	user1 := createRandomUser(t)

	// create a new user using the same details and verify if the email is the same: it should return an error
	user2, err := testQuery.CreateUser(context.Background(), db.CreateUserParams{
		Email:          user1.Email,
		HashedPassword: user1.HashedPassword,
	})

	assert.Error(t, err)
	assert.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	newPassword, err := utils.GenerateHashPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("unable to generate password:", err)
	}

	arg := db.UpdateUserPasswordParams{
		HashedPassword: newPassword,
		ID:             user.ID,
		UpdatedAt:      time.Now(),
	}

	newUser, err := testQuery.UpdateUserPassword(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, newUser.HashedPassword, arg.HashedPassword)
	assert.Equal(t, user.Email, newUser.Email)

	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
}

func TestGetUserByID(t *testing.T) {
	user := createRandomUser(t)

	newUser, err := testQuery.GetUserByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, newUser)

	assert.Equal(t, newUser.Email, user.Email)
	assert.Equal(t, newUser.HashedPassword, user.HashedPassword)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	newUser, err := testQuery.GetUserByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotEmpty(t, newUser)

	assert.Equal(t, newUser.Email, user.Email)
	assert.Equal(t, newUser.HashedPassword, user.HashedPassword)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQuery.DeleteUser(context.Background(), user.ID)

	assert.NoError(t, err)

	newUser, err := testQuery.GetUserByID(context.Background(), user.ID)

	assert.Error(t, err)
	assert.Empty(t, newUser)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := db.ListUsersParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := testQuery.ListUsers(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, len(users), 10)
}
