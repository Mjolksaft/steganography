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

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserInfo :one
SELECT 
    u.id,
    u.username,
    u.created_at,
    JSON_AGG(p.application_name) AS applications
FROM 
    users u
JOIN 
    passwords p ON p.user_id = u.id
WHERE 
    u.id = $1
GROUP BY 
    u.id, u.username;