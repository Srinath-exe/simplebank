package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Srinath-exe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.Amount, args.Amount)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.FromAccountID, args.FromAccountID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	getTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getTransfer)
	require.Equal(t, transfer.Amount, getTransfer.Amount)
	require.Equal(t, transfer.FromAccountID, getTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, getTransfer.ToAccountID)
	require.Equal(t, transfer.ID, getTransfer.ID)
	require.WithinDuration(t, transfer.CreatedAt, getTransfer.CreatedAt, time.Second)
}

func TestListTransfersFromAccount(t *testing.T) {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t)
	}

	for i := 0; i < 5; i++ {
		arg := CreateTransferParams{
			FromAccountID: account.ID,
			ToAccountID:   account2.ID,
			Amount:        util.RandomMoney(),
		}
		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// List transfers from account
	arg := ListTransfersFromAccountIdParams{
		FromAccountID: account.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfersFromAccountId(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account.ID, transfer.FromAccountID)
	}

	// // List transfers to account
	// toarg := ListTransfersToAccountIdParams{
	// 	ToAccountID: account2.ID,
	// 	Limit:       5,
	// 	Offset:      0,
	// }

	// transfers, err = testQueries.ListTransfersToAccountId(context.Background(), toarg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account2.ID, transfer.ToAccountID)
	}

}

func TestSeachTransfersByAccountOwner(t *testing.T) {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t)
	}

	for i := 0; i < 5; i++ {
		arg := CreateTransferParams{
			FromAccountID: account.ID,
			ToAccountID:   account2.ID,
			Amount:        util.RandomMoney(),
		}
		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// List transfers from account
	arg := SeachTransfersByAccountOwnerParams{
		SearchQuery: sql.NullString{String: account.Owner, Valid: true},
		Limit:       5,
		Offset:      0}

	transfers, err := testQueries.SeachTransfersByAccountOwner(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)

	}

}
