package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	n := 100
	amount := int64(10)

	processing := make(map[int64]struct{})
	type beforeAccounts struct {
		from, to Account
	}
	berforeAcntMap := make(map[int64]*beforeAccounts)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		account1 := getRandomAccount(t)
		account2 := getRandomAccount(t)
		id1, ok1 := processing[account1.ID]
		id2, ok2 := processing[account2.ID]
		if ok1 || ok2 {
			if ok1 {
				t.Logf("Already exist ID : %d\n", id1)

			} else if ok2 {
				t.Logf("Already exist ID : %d\n", id2)
			}
			continue
		}
		processing[account1.ID] = struct{}{}
		processing[account2.ID] = struct{}{}
		berforeAcntMap[account1.ID] = &beforeAccounts{account1, account2}
		txName := fmt.Sprintf("Tx %02d", i+1)
		t.Logf("%s - %d, %d", txName, account1.ID, account2.ID)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
		delete(processing, account1.ID)
		delete(processing, account2.ID)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		actnts := berforeAcntMap[result.Transfer.FromAccountID]
		account1, account2 := actnts.from, actnts.to

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check entries
		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		require.Equal(t, account1.ID, FromEntry.AccountID)
		require.Equal(t, -amount, FromEntry.Amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		// chefck the final updated balance
		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)

		require.Equal(t, account1.Balance-amount, updatedAccount1.Balance)
		require.Equal(t, account2.Balance+amount, updatedAccount2.Balance)
	}
}
