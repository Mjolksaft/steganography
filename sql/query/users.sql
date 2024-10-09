-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, hashed_password) 
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUsers :many
SELECT * FROM users;