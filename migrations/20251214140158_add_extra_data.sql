-- +goose up
ALTER TABLE users ADD COLUMN password text;
