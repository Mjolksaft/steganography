-- +goose Up
ALTER TABLE passwords
RENAME COLUMN application to application_name;

-- +goose Down
ALTER TABLE passwords
RENAME COLUMN application_name to application;
