-- +goose Up
-- SQLite does not support ENUM types.
-- Enum validation is handled via TEXT columns with CHECK constraints
-- in the respective table definitions.

-- +goose Down
