package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	require.NoError(t, CheckPassword(password, hashed))

	wrongPassword := RandomString(6)
	require.EqualError(t, bcrypt.ErrMismatchedHashAndPassword, CheckPassword(wrongPassword, hashed).Error())
}
