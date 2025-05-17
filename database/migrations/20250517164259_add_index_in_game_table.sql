-- +goose Up
-- +goose StatementBegin
-- add index in slug, title and description columns
CREATE INDEX idx_title ON games (title);
CREATE INDEX idx_description ON games (description);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_slug ON games;
DROP INDEX idx_title ON games;
DROP INDEX idx_description ON games;
-- +goose StatementEnd