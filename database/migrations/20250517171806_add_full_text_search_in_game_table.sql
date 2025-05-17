-- +goose Up
-- +goose StatementBegin
-- add full text search in game table
ALTER TABLE games
ADD COLUMN IF NOT EXISTS search_vector TSVECTOR;
-- add index in search_vector column
CREATE INDEX IF NOT EXISTS idx_search_vector ON games USING GIN (search_vector);
-- update search_vector column with existing data
UPDATE games
SET search_vector = to_tsvector(
        'english',
        coalesce(title, '') || ' ' || coalesce(description, '') || ' '
    );
-- create trigger to update search_vector column on insert or update
CREATE FUNCTION games_search_vector_update() RETURNS trigger AS $$ BEGIN NEW.search_vector := to_tsvector(
    'english',
    coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, '')
);
RETURN NEW;
END $$ LANGUAGE plpgsql;
CREATE TRIGGER games_search_vector_trigger BEFORE
INSERT
    OR
UPDATE ON games FOR EACH ROW EXECUTE FUNCTION games_search_vector_update();
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER games_search_vector_trigger ON games;
DROP FUNCTION games_search_vector_update();
ALTER TABLE games DROP COLUMN search_vector;
-- +goose StatementEnd