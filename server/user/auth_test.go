package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	userID := 1
	token, err := generateToken(uint32(userID), false)
	require.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateTokenWithInvalidSigningSecret(t *testing.T) {
	signingSecret = "hello world"
	userID := 1
	tokenString, err := generateToken(uint32(userID), true)
	require.NotNil(t, err)
	assert.EqualError(t, err, "key is of invalid type")
	assert.Empty(t, tokenString)
}

func TestValidateToken(t *testing.T) {
	signingSecret = []byte("hello world")
	userID := 1
	tokenString, err := generateToken(uint32(userID), false)
	require.Nil(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestValidateInvalidToken(t *testing.T) {
	authUserID, err := validateToken("hello one two three")
	require.NotNil(t, err)
	require.Equal(t, uint32(0), authUserID)
}

func TestValidateExpiredToken(t *testing.T) {
	signingSecret = []byte("hello world")
	userID := 1
	tokenString, err := generateToken(uint32(userID), true)
	require.Nil(t, err)
	assert.NotEmpty(t, tokenString)

	authUserID, err := validateToken(tokenString)
	require.NotNil(t, err)
	assert.EqualError(t, err, "expired token")
	require.Equal(t, uint32(0), authUserID)
}

func TestValidateTokenWithInvalidSigningSecret(t *testing.T) {
	signingSecret = []byte("hello world")
	userID := 1
	tokenString, err := generateToken(uint32(userID), false)
	require.Nil(t, err)
	assert.NotEmpty(t, tokenString)

	signingSecret = "hello world"

	authUserID, err := validateToken(tokenString)
	require.NotNil(t, err)
	assert.EqualError(t, err, "key is of invalid type")
	require.Equal(t, uint32(0), authUserID)
}
