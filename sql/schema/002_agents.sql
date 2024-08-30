-- +goose Up
ALTER TABLE agents
ADD COLUMN reset_datetime INT;

-- +goose Down
ALTER TABLE agents
DROP COLUMN reset_datetime;
