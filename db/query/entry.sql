-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
    ) VALUES (
    $1,
    $2
    ) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries WHERE id = $1 LIMIT 1;


-- name: ListEntryFromAccountId :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;

-- name: SeachEntriesByAccountOwner :many
SELECT e.*
FROM entries e
INNER JOIN accounts a ON e.account_id = a.id
WHERE a.owner ILIKE '%' || sqlc.arg(search_query) || '%'
AND e.created_at >= sqlc.arg(start_date) AND e.created_at <= sqlc.arg(end_date)
AND e.amount >= sqlc.arg(min_amount) AND e.amount <= sqlc.arg(max_amount)
ORDER BY
CASE WHEN  sqlc.arg(field) = 'amount' AND  sqlc.arg(order_by) = 'ASC' THEN e.amount END  ASC,
CASE WHEN  sqlc.arg(field) = 'amount' AND  sqlc.arg(order_by) = 'DESC' THEN e.amount END DESC,
CASE WHEN  sqlc.arg(field) = 'created_at' AND  sqlc.arg(order_by) = 'ASC' THEN e.created_at END  ASC,
CASE WHEN  sqlc.arg(field) = 'created_at' AND  sqlc.arg(order_by) = 'DESC' THEN e.created_at END DESC,
CASE WHEN  sqlc.arg(field) = 'id' AND  sqlc.arg(order_by) = 'ASC' THEN e.id END  ASC,
CASE WHEN  sqlc.arg(field) = 'id' AND  sqlc.arg(order_by) = 'DESC' THEN e.id END DESC
LIMIT $1
OFFSET $2;

