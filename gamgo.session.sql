-- Example: Filter by title and platforms (case-insensitive, partial match)
-- EXPLAIN 
SELECT *
FROM games
WHERE (
        '' IS NULL
        OR search_vector @@ plainto_tsquery('')
    )
    OR (
        'Web' IS NULL
        OR platforms @> ARRAY ['Web']
    );
---
-- Example: Filter by title (using search_vector) and platforms (case-insensitive, partial match)
EXPLAIN
SELECT *
FROM games
WHERE (
        'auto' IS NULL
        OR search_vector @@ plainto_tsquery('auto')
    )
    AND (
        '' IS NULL
        OR platforms @> ARRAY ['']
    );
---
SELECT column_name,
    data_type
FROM information_schema.columns
WHERE table_name = 'games';
---
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS goose_db_version;