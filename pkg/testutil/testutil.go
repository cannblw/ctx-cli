package testutil

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/cannblw/ctx-cli/migrations"
	"github.com/cannblw/ctx-cli/pkg/store"
)

func NewTestDB(t *testing.T) *store.Store {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	sqldb, err := sql.Open(store.DriverName, dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { sqldb.Close() })

	goose.SetDialect(store.GooseDialect)
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations.FS)
	require.NoError(t, goose.Up(sqldb, "."))

	_, err = sqldb.Exec(store.ForeignKeysPragma)
	require.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() { db.Close() })
	return store.NewStoreFromDB(db)
}
