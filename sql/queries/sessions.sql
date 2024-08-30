-- name: CreateSession :exec
INSERT INTO sessions (id, expires_at, user_id, agent_id)
VALUES ($1, $2, $3, $4);

-- name: GetSessionById :one
SELECT * FROM sessions 
WHERE id = $1 and expires_at > NOW();

-- name: SetAgentForSession :exec
UPDATE sessions 
SET agent_id = $1
WHERE id = $2;
