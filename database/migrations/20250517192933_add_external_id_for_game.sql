-- +goose Up
-- +goose StatementBegin
ALTER TABLE games
ADD COLUMN externalId VARCHAR(255) DEFAULT NULL,
ADD COLUMN externalSource VARCHAR(255) DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE games
DROP COLUMN externalId,
DROP COLUMN externalSource;

-- +goose StatementEnd