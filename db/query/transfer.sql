-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
    ) VALUES (  
    $1,$2,$3
    ) RETURNING *;

-- name: GetTransfer :one    
SELECT * FROM transfers WHERE id = $1 LIMIT 1;


-- name: ListTransfersFromAccountId :many
SELECT t.*, 
json_build_object('owner', a1.owner, 'balance', a1.balance) AS from_account,
json_build_object('owner', a2.owner, 'balance', a2.balance) AS to_account
FROM transfers t
INNER JOIN accounts a1 ON t.from_account_id = a1.id
INNER JOIN accounts a2 ON t.to_account_id = a2.id
WHERE t.from_account_id = $1
OR t.to_account_id = $1
ORDER BY t.created_at
LIMIT $2
OFFSET $3;

-- name: ListTransfersToAccountId :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: SeachTransfersByAccountOwner :many
SELECT t.* , 
json_build_object('owner', a1.owner, 'balance', a1.balance) AS from_account,
json_build_object('owner', a2.owner, 'balance', a2.balance) AS to_account
FROM transfers t
INNER JOIN accounts a1 ON t.from_account_id = a1.id
INNER JOIN accounts a2 ON t.to_account_id = a2.id
WHERE a1.owner ILIKE '%' || sqlc.arg(search_query) || '%'
OR a2.owner ILIKE '%' || sqlc.arg(search_query) || '%'
LIMIT $1
OFFSET $2;


