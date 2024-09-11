package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Srinath-exe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Tranfer {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := CreateTranferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomMoney(),
	}

	transfer,err := testQueries.CreateTranfer(context.Background(),args)
	require.NoError(t, err)
	require.NotEmpty(t,transfer)
	require.Equal(t,transfer.Amount,args.Amount)
	require.Equal(t,transfer.ToAccountID,args.ToAccountID)
	require.Equal(t,transfer.FromAccountID,args.FromAccountID)

	return transfer
}

func TestCreateTransfer(t *testing.T){
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T){
	transfer := createRandomTransfer(t)

	getTransfer,err := testQueries.GetTransfer(context.Background(),transfer.ID)

	require.NoError(t,err)
	require.NotEmpty(t,getTransfer)
	require.Equal(t,transfer.Amount,getTransfer.Amount)
	require.Equal(t,transfer.FromAccountID,getTransfer.FromAccountID)
	require.Equal(t,transfer.ToAccountID,getTransfer.ToAccountID)
	require.Equal(t,transfer.ID,getTransfer.ID)
	require.WithinDuration(t, transfer.CreatedAt,getTransfer.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T){
  transfer :=	createRandomTransfer(t)
	args := UpdateTransferParams{
		ID: transfer.ID,
		Amount: util.RandomMoney(),
	}

	getTransfer,err := testQueries.UpdateTransfer(context.Background(),args)
	
	require.NoError(t,err)
	require.NotEmpty(t,getTransfer)
	require.Equal(t,args.Amount,getTransfer.Amount)
	require.Equal(t,transfer.FromAccountID,getTransfer.FromAccountID)
	require.Equal(t,transfer.ToAccountID,getTransfer.ToAccountID)
	require.Equal(t,transfer.ID,getTransfer.ID)
	require.WithinDuration(t, transfer.CreatedAt,getTransfer.CreatedAt, time.Second)
	

}

func TestDeleteTransfer(t *testing.T){
	tranfer := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(),tranfer.ID)
	require.NoError(t,err)

	getTransfer ,err := testQueries.GetTransfer(context.Background(),tranfer.ID)

	require.Error(t,err)
	require.Error(t,err,sql.ErrNoRows.Error())
	require.Empty(t, getTransfer)

}

func TestListTransfer(t *testing.T){
	
	for i := 0; i < 10 ; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit: 5,
		Offset: 3,
	}

	tranfers ,err := testQueries.ListTransfers(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,tranfers)
	require.Len(t,tranfers,5)

	for _,transfer := range tranfers {
		require.NotEmpty(t,transfer)
	}

}

func TestListTransfersFromAccount(t *testing.T){
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	
	for i := 0; i < 5 ; i++ {
		createRandomTransfer(t)
	}


for i := 0; i < 5; i++ {
		arg := CreateTranferParams{
			FromAccountID: account.ID,
			ToAccountID: account2.ID,
			Amount:    util.RandomMoney(),
		}
		_, err := testQueries.CreateTranfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// List transfers from account
	arg := ListTransfersFromAccountIdParams{
		FromAccountID: account.ID,
		Limit: 5,
		Offset: 0,
	}

	transfers ,err := testQueries.ListTransfersFromAccountId(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,transfers)

	require.Len(t,transfers,5)

	for _,transfer := range transfers {
		require.NotEmpty(t,transfer)
		require.Equal(t,account.ID,transfer.FromAccountID)
	}

	// List transfers to account
	toarg := ListTransfersToAccountIdParams{
		ToAccountID: account2.ID,
		Limit: 5,
		Offset: 0,
	}

	transfers ,err = testQueries.ListTransfersToAccountId(context.Background(),toarg)

	require.NoError(t,err)
	require.NotEmpty(t,transfers)

	require.Len(t,transfers,5)

	for _,transfer := range transfers {
		require.NotEmpty(t,transfer)
		require.Equal(t,account2.ID,transfer.ToAccountID)
	}

}