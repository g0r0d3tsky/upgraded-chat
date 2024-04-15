-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE messages
ADD COLUMN time TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE messages
DROP COLUMN time;
-- +goose StatementEnd