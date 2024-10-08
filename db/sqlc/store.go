package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	DeleteUserWithAccountsTx(ctx context.Context, username string) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rberr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction

type TransferTxParams struct {
	FromAccID int64 `json:"from_account_id"`
	ToAccID   int64 `json:"to_account_id"`
	Amount    int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccID,
			ToAccountID:   arg.ToAccID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccID < arg.ToAccID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccID, -arg.Amount, arg.ToAccID, arg.Amount)
			if err != nil {
				return err
			}

		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccID, arg.Amount, arg.FromAccID, -arg.Amount)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

func AddMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	if err != nil {
		return
	}

	return
}

func (store *SQLStore) DeleteUserWithAccountsTx(ctx context.Context, username string) error {

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.GetUser(ctx, username)

		if err != nil {
			return err
		}

		accounts, err := q.ListAccounts(ctx, ListAccountsParams{Owner: username})

		fmt.Println(accounts)

		if err != nil {
			return err
		}

		for _, account := range accounts {
			err = q.DeleteAccount(ctx, account.ID)
			fmt.Println(account.ID)

			if err != nil {
				return err
			}
		}

		fmt.Println(username)
		err = q.DeleteUser(ctx, username)

		if err != nil {
			return err
		}

		return nil
	})

	return err
}
