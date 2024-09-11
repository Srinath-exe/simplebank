-- name: CreateTranfer :one
INSERT INTO tranfers (
    from_account_id,
    to_account_id,
    amount
    ) VALUES (  
    $1,$2,$3
    ) RETURNING *;

-- name: GetTransfer :one    
SELECT * FROM tranfers WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM tranfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListTransfersFromAccountId :many
SELECT * FROM tranfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListTransfersToAccountId :many
SELECT * FROM tranfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateTransfer :one
UPDATE tranfers
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM tranfers WHERE id = $1;


