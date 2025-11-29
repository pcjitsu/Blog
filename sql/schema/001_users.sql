-- +goose Up
-- This section runs when you migrate UP (create tables).
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
-- This section runs when you migrate DOWN (rollback/delete tables).
DROP TABLE users;