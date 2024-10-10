-- +goose Up
ALTER TABLE passwords
ADD CONSTRAINT unique_application UNIQUE (application);

-- +goose Down
ALTER TABLE passwords
DROP CONSTRAINT unique_application;
