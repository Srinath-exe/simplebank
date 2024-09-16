package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Srinath-exe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createAccount := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), createAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, createAccount.ID, account.ID)
	require.Equal(t, createAccount.Owner, account.Owner)
	require.Equal(t, createAccount.Balance, account.Balance)
	require.Equal(t, createAccount.Currency, account.Currency)
	require.WithinDuration(t, createAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createAccount := createRandomAccount(t)
	arg := AddAccountBalanceParams{
		ID:     createAccount.ID,
		Amount: util.RandomMoney(),
	}
	account, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	newBal := createAccount.Balance + arg.Amount
	require.NotEmpty(t, account)
	require.Equal(t, createAccount.ID, account.ID)
	require.Equal(t, createAccount.Owner, account.Owner)
	require.Equal(t, newBal, account.Balance)
	require.Equal(t, createAccount.Currency, account.Currency)
	require.WithinDuration(t, createAccount.CreatedAt, account.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	createAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), createAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), createAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestSearchAccounts(t *testing.T) {
	var account Account
	for i := 0; i < 10; i++ {
		account = createRandomAccount(t)
	}

	arg := SearchAccountsParams{
		Column1: sql.NullString{String: account.Owner, Valid: true},
		Limit:   5,
		Offset:  0,
	}

	accounts, err := testQueries.SearchAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 1)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, arg.Column1.String, account.Owner)
	}
}
