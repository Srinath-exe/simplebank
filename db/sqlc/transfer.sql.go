// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transfer.sql

package db

import (
	"context"
)

const createTranfer = `-- name: CreateTranfer :one
INSERT INTO tranfers (
    from_account_id,
    to_account_id,
    amount
    ) VALUES (  
    $1,$2,$3
    ) RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTranferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *Queries) CreateTranfer(ctx context.Context, arg CreateTranferParams) (Tranfer, error) {
	row := q.db.QueryRowContext(ctx, createTranfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Tranfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM tranfers WHERE id = $1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM tranfers WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Tranfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Tranfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM tranfers
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListTransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Tranfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Tranfer{}
	for rows.Next() {
		var i Tranfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfersFromAccountId = `-- name: ListTransfersFromAccountId :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM tranfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTransfersFromAccountIdParams struct {
	FromAccountID int64 `json:"from_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) ListTransfersFromAccountId(ctx context.Context, arg ListTransfersFromAccountIdParams) ([]Tranfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfersFromAccountId, arg.FromAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Tranfer{}
	for rows.Next() {
		var i Tranfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfersToAccountId = `-- name: ListTransfersToAccountId :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM tranfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTransfersToAccountIdParams struct {
	ToAccountID int64 `json:"to_account_id"`
	Limit       int32 `json:"limit"`
	Offset      int32 `json:"offset"`
}

func (q *Queries) ListTransfersToAccountId(ctx context.Context, arg ListTransfersToAccountIdParams) ([]Tranfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfersToAccountId, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Tranfer{}
	for rows.Next() {
		var i Tranfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransfer = `-- name: UpdateTransfer :one
UPDATE tranfers
SET amount = $2
WHERE id = $1
RETURNING id, from_account_id, to_account_id, amount, created_at
`

type UpdateTransferParams struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (q *Queries) UpdateTransfer(ctx context.Context, arg UpdateTransferParams) (Tranfer, error) {
	row := q.db.QueryRowContext(ctx, updateTransfer, arg.ID, arg.Amount)
	var i Tranfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
