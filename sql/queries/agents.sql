-- name: CreateAgent :exec
INSERT INTO agents (id, name, token, user_id)
VALUES ($1, $2, $3, $4);

-- name: GetAgentsByUserId :many
SELECT id, name FROM agents WHERE user_id = $1;

-- name: GetAgentTokenById :one
SELECT token FROM agents WHERE id = $1;
