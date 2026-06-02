package repo_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/cannblw/ctx-cli/pkg/models"
)

func findMigrationsDir() string {
	dirs := []string{"migrations", "../migrations", "../../migrations", "../../../migrations"}
	for _, d := range dirs {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			return d
		}
	}
	return "migrations"
}

func newTestBunDB(t *testing.T) *bun.DB {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	sqldb, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { sqldb.Close() })

	goose.SetDialect("sqlite3")
	goose.SetLogger(goose.NopLogger())
	require.NoError(t, goose.Up(sqldb, findMigrationsDir()))

	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() { db.Close() })
	return db
}

func TestMigration_CreatesDefaultStates(t *testing.T) {
	db := newTestBunDB(t)

	var states []models.State
	err := db.NewSelect().Model(&states).OrderExpr("position ASC").Scan(context.Background())
	require.NoError(t, err)
	require.Len(t, states, 4)
	assert.Equal(t, "todo", states[0].Name)
	assert.Equal(t, "in-progress", states[1].Name)
	assert.Equal(t, "review", states[2].Name)
	assert.Equal(t, "done", states[3].Name)
}

func TestMigration_Idempotent(t *testing.T) {
	db := newTestBunDB(t)

	sqldb := db.DB
	goose.SetDialect("sqlite3")
	goose.SetLogger(goose.NopLogger())
	require.NoError(t, goose.Up(sqldb, findMigrationsDir()))

	count, err := db.NewSelect().Model((*models.State)(nil)).Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 4, count, "migration should be idempotent")
}
