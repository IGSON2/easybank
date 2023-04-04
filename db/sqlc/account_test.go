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
	user := TestGetRandomUser(t)
	currencies, err := testQueries.GetCurrency(context.Background(), owner)
	require.NoError(t, err)

	arg := CreateAccountParams{
		Owner:    user.Username, //user 테이블과 FK 제약조건이 걸려있기 때문에 account 먼저 생성될 수 없음.
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(currencies),
	}
	require.NotEqual(t, "", arg.Currency)

	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NotNil(t, result)
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
	for i := 0; i < 8; i++ {
		testCreateAccount(t)
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
