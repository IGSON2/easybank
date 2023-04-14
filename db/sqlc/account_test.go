package db

import (
	"context"
	"database/sql"
	"easybank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	owner    = "ginber"
	password = "password"
)

func testCreateAccount(t *testing.T) Account {
	user := testGetRandomUser(t)
	currencies, err := testQueries.GetCurrency(context.Background(), user.Username)
	require.NoError(t, err)

	if len(currencies) == len(util.Currencies) {
		account, err := testQueries.GetAccount(context.Background(), GetAccountParams{
			Owner:    user.Username,
			Currency: currencies[util.RandomInt(0, int64(len(currencies)-1))],
		})
		require.NoError(t, err)
		return account
	}

	arg := CreateAccountParams{
		Owner:    user.Username, //user 테이블과 FK 제약조건이 걸려있기 때문에 account 먼저 생성될 수 없음.
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(currencies),
	}

	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NotNil(t, result, err)
	require.NoError(t, err)

	getArg := GetAccountParams{
		Owner:    arg.Owner,
		Currency: arg.Currency,
	}

	account, err := testQueries.GetAccount(context.Background(), getArg)
	require.NoError(t, err)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	testCreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	getFirstAccount(t)
}

func TestUpdateAccount(t *testing.T) {
	account := getFirstAccount(t)
	randomMoney := util.RandomBalance()
	result, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		Balance: account.Balance + randomMoney,
		ID:      account.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	getArg := GetAccountParams{
		Owner:    account.Owner,
		Currency: account.Currency,
	}

	account2, err := testQueries.GetAccount(context.Background(), getArg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.WithinDuration(t, account.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account.Balance+randomMoney, account2.Balance)
}

func TestDeleteAcc(t *testing.T) {
	account := getFirstAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	getArg := GetAccountParams{
		Owner:    account.Owner,
		Currency: account.Currency,
	}

	account2, err := testQueries.GetAccount(context.Background(), getArg)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = testCreateAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}

}

func getFirstAccount(t *testing.T) Account {
	currencies, err := testQueries.GetCurrency(context.Background(), owner)
	require.NoError(t, err)
	require.NotNil(t, currencies)

	getArg := GetAccountParams{
		Owner:    owner,
		Currency: currencies[0],
	}

	account, err := testQueries.GetAccount(context.Background(), getArg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, getArg.Owner, account.Owner)
	require.Equal(t, getArg.Currency, account.Currency)

	return account
}
