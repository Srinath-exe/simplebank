package db

import (
	"context"
	"database/sql"
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

func TestDeleteUser(t *testing.T) {
	createUser := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), createUser.Username)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), createUser.Username)
	require.Error(t, err)
	require.Empty(t, user)
}

func TestUpdatePassword(t *testing.T) {
	user := createRandomUser(t)
	newpsw := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newpsw)
	require.NoError(t, err)

	arg := UpdatePasswordParams{
		Username:       user.Username,
		HashedPassword: hashedPassword,
	}
	err = testQueries.UpdatePassword(context.Background(), arg)
	require.NoError(t, err)

	user, err = testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.NoError(t, util.CheckPasswordHash(newpsw, user.HashedPassword))
}

func TestSearchUsers(t *testing.T) {
	var user User
	for i := 0; i < 10; i++ {
		user = createRandomUser(t)
	}

	arg := SearchUsersParams{
		Column1: sql.NullString{String: user.Username, Valid: true},
		Limit:   5,
		Offset:  0,
	}
	users, err := testQueries.SearchUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 1)

	for _, user := range users {
		require.NotEmpty(t, user)
		require.Contains(t, user.Username, arg.Column1.String)
	}
}

func TestGetUsersParams(t *testing.T) {
	var users []User

	for i := 0; i < 10; i++ {
		users = append(users, createRandomUser(t))
	}

	arg := GetUsersParams{
		Limit:  5,
		Offset: 0,
		Usernames: []string{
			users[5].Username,
			users[6].Username,
			users[7].Username,
		},
	}

	selectedUsers, err := testQueries.GetUsers(context.Background(), arg)

	require.NoError(t, err)

	require.Len(t, selectedUsers, 3)

	for i, user := range selectedUsers {
		require.Contains(t, arg.Usernames, user.Username)
		require.Contains(t, users, selectedUsers[i])
	}

}
