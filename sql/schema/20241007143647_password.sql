-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    hashed_password TEXT,
    username TEXT
);

CREATE TABLE passwords (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    hashed_password TEXT,
    application TEXT,
    user_id UUID REFERENCES users(id)
);


-- +goose Down
DROP TABLE passwords;
DROP TABLE users;