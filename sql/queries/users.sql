-- name: CreateUser :one
-- The comment above tells sqlc to generate a function named 'CreateUser' that returns one row.
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1, -- Corresponds to arg.ID in Go
    $2, -- Corresponds to arg.CreatedAt in Go
    $3, -- Corresponds to arg.UpdatedAt in Go
    $4  -- Corresponds to arg.Name in Go
)
RETURNING *; -- Return the full created row so we have the data back in Go

-- name: GetUser :one
-- Generate a function 'GetUser' that selects a single user by name.
SELECT id, name, created_at, updated_at from users WHERE name = $1;