package db

import (
	"context"
	"samplebank/util"

	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	checkAccount := CreateAccountParams{
		Username: util.RandomUsername(),
		Balance:  util.RandomBalance(),
		Currency: "USD",
	}

	testAccount, err := testQueries.CreateAccount(context.Background(), checkAccount)
	require.NoError(t, err)
	require.NotEmpty(t, testAccount)
	return testAccount
}
func TestCreateAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	require.NotEmpty(t, testAccount)

	require.NotZero(t, testAccount.ID)
	require.NotZero(t, testAccount.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, account2.Balance)
}
