// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: passwords.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createPassword = `-- name: CreatePassword :exec
INSERT INTO passwords (id, created_at, updated_at, hashed_password, application, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
`

type CreatePasswordParams struct {
	HashedPassword string
	Application    string
	UserID         uuid.UUID
}

func (q *Queries) CreatePassword(ctx context.Context, arg CreatePasswordParams) error {
	_, err := q.db.ExecContext(ctx, createPassword, arg.HashedPassword, arg.Application, arg.UserID)
	return err
}
