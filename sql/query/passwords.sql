-- name: CreatePassword :exec
INSERT INTO passwords (id, created_at, updated_at, hashed_password, application_name, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
);

-- name: GetPasswords :many
SELECT * FROM passwords
WHERE user_id = $1;


-- name: GetPassword :one
SELECT * FROM passwords
WHERE application_name = $1 AND user_id = $2;
