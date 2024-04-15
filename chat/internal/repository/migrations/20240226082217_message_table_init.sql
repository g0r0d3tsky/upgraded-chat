-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS messages (
    id uuid PRIMARY KEY,
    content TEXT NOT NULL,
    user_nickname VARCHAR(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd