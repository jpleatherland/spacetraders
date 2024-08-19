-- name: CreateUser :one
INSERT INTO users (id, name, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1; 
