package db

import (
	"context"
	"easybank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	result, err := testQueries.CreateUser(context.Background(), arg)
	require.NotNil(t, result)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), arg.Username)
	require.NoError(t, err)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func testGetRandomUser(t *testing.T) User {
	user, err := testQueries.GetRandomUser(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}
