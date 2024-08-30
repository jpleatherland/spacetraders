-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE agents (
    id UUID primary key,
    name TEXT NOT NULL,
    token TEXT NOT NULL,
    user_id UUID NOT NULL,
    reset_datetime INT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE sessions (
    id TEXT primary key,
    expires_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    agent_id UUID,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);

-- +goose Down
DROP TABLE sessions;
DROP TABLE agents;
DROP TABLE users;
