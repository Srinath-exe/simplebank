package db

import (
	"context"
	"testing"

	"github.com/Srinath-exe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T ) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	return entry
}


func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntryParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, arg.Amount, entry2.Amount)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)
}

func TestListEntriesFromAccount(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomEntry(t)
	}
	for i := 0; i < 5; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomMoney(),
		}
		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListEntryFromAccountIdParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.ListEntryFromAccountId(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}

func TestListEntriesFromAccountInvalid(t *testing.T) {
	arg := ListEntryFromAccountIdParams{
		AccountID: 0,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.ListEntryFromAccountId(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, entries)
}

