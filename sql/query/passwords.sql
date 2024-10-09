-- name: CreateOne :one
INSERT INTO password (id, created_at, hashed_password, application)
VALUES (
    gen_random_uuid(),
    NOW(),
    $1,
    $2
)
RETURNING *;