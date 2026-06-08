-- +goose Up
UPDATE items SET context_id = NULL WHERE context_id IN (SELECT id FROM contexts WHERE name = 'global');
DELETE FROM contexts WHERE name = 'global';

-- +goose Down
INSERT OR IGNORE INTO contexts (name, description) VALUES ('global', 'Default global context');
