-- +goose Up
INSERT OR IGNORE INTO contexts (name, description) VALUES ('global', 'Default global context');

-- +goose Down
DELETE FROM contexts WHERE name = 'global';
