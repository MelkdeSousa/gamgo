-- +goose Up
-- +goose StatementBegin
-- setup uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- create games table
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    platforms TEXT [] NOT NULL,
    releaseDate DATE NOT NULL,
    rating INTEGER NOT NULL,
    coverImage TEXT NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE games;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd