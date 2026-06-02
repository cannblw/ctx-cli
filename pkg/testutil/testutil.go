package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	_ "github.com/cannblw/ctx-cli/migrations"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func NewTestDB(t *testing.T) *bun.DB {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	sqldb, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { sqldb.Close() })

	goose.SetDialect("sqlite3")
	migrationsDir := findMigrationsDir()
	require.NoError(t, goose.Up(sqldb, migrationsDir))

	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() { db.Close() })
	return db
}

func findMigrationsDir() string {
	dirs := []string{"migrations", "../migrations", "../../migrations", "../../../migrations"}
	for _, d := range dirs {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			abspath, _ := filepath.Abs(d)
			return abspath
		}
	}
	return "migrations"
}
