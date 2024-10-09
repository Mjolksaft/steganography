-- +goose Up
CREATE TABLE password (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    hashed_password TEXT,
    application TEXT
);

-- +goose Down
DROP TABLE password;