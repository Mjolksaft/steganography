// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: passwords.sql

package database

import (
	"context"
	"database/sql"
)

const createOne = `-- name: CreateOne :one
INSERT INTO password (id, created_at, hashed_password, application)
VALUES (
    gen_random_uuid(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, hashed_password, application
`

type CreateOneParams struct {
	HashedPassword sql.NullString
	Application    sql.NullString
}

func (q *Queries) CreateOne(ctx context.Context, arg CreateOneParams) (Password, error) {
	row := q.db.QueryRowContext(ctx, createOne, arg.HashedPassword, arg.Application)
	var i Password
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.Application,
	)
	return i, err
}
