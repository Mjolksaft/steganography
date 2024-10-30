-- +goose Up
ALTER TABLE passwords
ALTER COLUMN Created_At SET NOT NULL,
ALTER COLUMN Updated_At SET NOT NULL;

ALTER TABLE users
ALTER COLUMN Created_At SET NOT NULL,
ALTER COLUMN Updated_At SET NOT NULL,
ALTER COLUMN is_admin SET NOT NULL;

-- +goose Down
ALTER TABLE passwords
ALTER COLUMN Created_At DROP NOT NULL,
ALTER COLUMN Updated_At DROP NOT NULL;

ALTER TABLE users
ALTER COLUMN Created_At DROP NOT NULL,
ALTER COLUMN Updated_At DROP NOT NULL,
ALTER COLUMN is_admin DROP NOT NULL; 
