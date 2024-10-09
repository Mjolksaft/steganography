-- +goose Up
ALTER TABLE users
ADD is_admin BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_admin;
