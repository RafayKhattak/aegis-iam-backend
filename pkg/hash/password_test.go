package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPasswordProducesDifferentHashesForSameInput(t *testing.T) {
	password := "enterprise_secure_123"

	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)

	require.NotEqual(t, hash1, hash2)
	require.NoError(t, CheckPassword(password, hash1))
	require.NoError(t, CheckPassword(password, hash2))
}

func TestCheckPassword(t *testing.T) {
	password := "enterprise_secure_123"
	wrongPassword := "not_the_right_password"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	require.NoError(t, CheckPassword(password, hashedPassword))
	require.Error(t, CheckPassword(wrongPassword, hashedPassword))
}
