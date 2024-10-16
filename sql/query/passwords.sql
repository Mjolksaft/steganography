-- name: CreatePassword :exec
INSERT INTO passwords (id, created_at, updated_at, hashed_password, application, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
);
