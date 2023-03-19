package db

import (
	"context"
	"database/sql"
	"easybank/util"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username, //user 테이블과 FK 제약조건이 걸려있기 때문에 account 먼저 생성될 수 없음.
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NotNil(t, result)
	require.NoError(t, err)

	lastID, err := result.LastInsertId()
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), lastID)
	require.NoError(t, err)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func getRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NotNil(t, result)
	require.NoError(t, err)

	lastID, err := result.LastInsertId()
	require.NoError(t, err)

	ranNum := util.RandomInt(1, lastID-1)

	account, err := testQueries.GetAccount(context.Background(), ranNum)
	require.NoError(t, err, fmt.Sprintf("Get random account Error. ID : %d", ranNum))

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	randomMoney := util.RandomBalance()
	result, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		Balance: account1.Balance + randomMoney,
		ID:      account1.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Balance+randomMoney, account2.Balance)
}

func TestDeleteAcc(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 8; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 3, // 첫 3개 레코드를 건너뛰고 5개를 반환
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, int(arg.Limit))

	for _, a := range accounts {
		require.NotEmpty(t, a)
	}

}
