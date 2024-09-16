-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
    ) VALUES (
    $1,
    $2,
    $3,
    $4
    ) RETURNING *;


-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE username = ANY(sqlc.arg(usernames)::text[])
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdatePassword :exec
UPDATE users
SET hashed_password = $2
WHERE username = $1;

-- name: SearchUsers :many
SELECT * FROM users
WHERE username ILIKE '%' || $1 || '%'
ORDER BY username
LIMIT $2
OFFSET $3;

-- name: DeleteUser :exec
DELETE FROM users WHERE username = $1;

