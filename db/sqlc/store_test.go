package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransactions(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)
	transaction, err := testQueries.CreateTransaction(context.Background(), CreateTransactionParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	fromRecord, err := testQueries.CreateRecord(context.Background(), CreateRecordParams{
		AccountID: account1.ID,
		Amount:    -amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, fromRecord)

	toRecord, err := testQueries.CreateRecord(context.Background(), CreateRecordParams{
		AccountID: account2.ID,
		Amount:    amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, toRecord)
}
func TestCreateTransaction(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 10
	amount := int64(10)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.transferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transaction := result.Transaction
		require.NotEmpty(t, transaction)
		require.Equal(t, account1.ID, transaction.FromAccountID)
		require.Equal(t, account2.ID, transaction.ToAccountID)
		require.Equal(t, amount, transaction.Amount)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.CreatedAt)

		_, err = store.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		// check record
		fromEntry := result.FromRecord
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetRecord(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToRecord
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetRecord(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		//check diff balance
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

	fmt.Println(">> After: ", account1.Balance, account2.Balance)
}

func TestCreateTransactionDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	errs := make(chan error)

	n := 10
	amount := int64(10)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.transferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1, updateAccount1.Balance)
	require.Equal(t, account2, updateAccount2.Balance)

	fmt.Println(">> After: ", account1.Balance, account2.Balance)
}
