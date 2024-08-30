-- name: CreateAgent :exec
INSERT INTO agents (id, name, token, user_id, reset_datetime)
VALUES ($1, $2, $3, $4, $5);

-- name: GetAgentsByUserId :many
SELECT id, name, reset_datetime FROM agents WHERE user_id = $1;

-- name: GetAgentTokenById :one
SELECT token FROM agents WHERE id = $1;

-- name: DeleteAgentById :exec
DELETE FROM agents WHERE id = $1;
