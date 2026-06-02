package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up001, Down001)
}

func Up001(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS states (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL UNIQUE,
			position   INTEGER NOT NULL,
			orphaned   INTEGER NOT NULL DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS contexts (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT    NOT NULL UNIQUE,
			description TEXT    NOT NULL DEFAULT '',
			created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS items (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			slug        TEXT    NOT NULL UNIQUE,
			context_id  INTEGER REFERENCES contexts(id) ON DELETE CASCADE,
			type        TEXT    NOT NULL,
			value       TEXT    NOT NULL,
			state_id    INTEGER REFERENCES states(id),
			notes       TEXT    NOT NULL DEFAULT '',
			created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
		);
		INSERT OR IGNORE INTO states (name, position, orphaned) VALUES
			('todo',        0, 0),
			('in-progress', 1, 0),
			('review',      2, 0),
			('done',        3, 0);
	`)
	return err
}

func Down001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS items; DROP TABLE IF EXISTS contexts; DROP TABLE IF EXISTS states;")
	return err
}
