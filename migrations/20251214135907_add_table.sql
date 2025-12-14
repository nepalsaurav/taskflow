-- +goose up
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    username text
);
