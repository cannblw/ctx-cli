package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/cannblw/ctx-cli/migrations"
)

const (
	DriverName        = "sqlite"
	GooseDialect      = "sqlite3"
	ForeignKeysPragma = "PRAGMA foreign_keys = ON"
)

type Store struct {
	db *bun.DB
}

func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create db dir: %w", err)
	}

	sqldb, err := sql.Open(DriverName, dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	goose.SetDialect(GooseDialect)
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations.FS)
	if err := goose.Up(sqldb, "."); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}

	if _, err := sqldb.Exec(ForeignKeysPragma); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("could not enable foreign keys: %w", err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	return &Store{db: db}, nil
}

func NewStoreFromDB(db *bun.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Close() error {
	return s.db.Close()
}
