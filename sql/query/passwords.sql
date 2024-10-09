-- name: CreatePassword :one
INSERT INTO passwords (id, created_at, updated_at, hashed_password, application, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetPassword :one
SELECT * FROM passwords
WHERE application = $1 AND user_id = $2;

-- name: GetPasswords :many
SELECT * FROM passwords;
