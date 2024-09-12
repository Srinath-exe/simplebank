package db

import (
	"context"
	"testing"
	"time"

	"github.com/Srinath-exe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createUser := createRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), createUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, createUser.Username, user.Username)
	require.Equal(t, createUser.HashedPassword, user.HashedPassword)
	require.Equal(t, createUser.Email, user.Email)
	require.Equal(t, createUser.FullName, user.FullName)
	require.WithinDuration(t, createUser.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, createUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
}
