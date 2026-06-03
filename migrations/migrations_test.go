package migrations_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/cannblw/ctx-cli/migrations"
	"github.com/cannblw/ctx-cli/pkg/models"
	"github.com/cannblw/ctx-cli/pkg/store"
)

func newTestBunDB(t *testing.T) *bun.DB {
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
	return db
}

func TestMigration_CreatesDefaultStates(t *testing.T) {
	db := newTestBunDB(t)

	var states []models.State
	err := db.NewSelect().Model(&states).OrderExpr("position ASC").Scan(context.Background())
	require.NoError(t, err)
	require.Len(t, states, 4)
	assert.Equal(t, "todo", states[0].Name)
	assert.Equal(t, 0, states[0].Position)
	assert.False(t, states[0].Orphaned)
	assert.Equal(t, "in-progress", states[1].Name)
	assert.Equal(t, 1, states[1].Position)
	assert.Equal(t, "review", states[2].Name)
	assert.Equal(t, 2, states[2].Position)
	assert.Equal(t, "done", states[3].Name)
	assert.Equal(t, 3, states[3].Position)
}

func TestMigration_CreatesContextsTable(t *testing.T) {
	db := newTestBunDB(t)

	c := &models.Context{
		Name:        "ctx-test",
		Description: "verify table exists",
	}
	_, err := db.NewInsert().Model(c).Exec(context.Background())
	require.NoError(t, err)
	assert.NotZero(t, c.ID)
	assert.False(t, c.CreatedAt.IsZero())
}

func TestMigration_CreatesItemsTable(t *testing.T) {
	db := newTestBunDB(t)

	c := &models.Context{Name: "items-test"}
	_, err := db.NewInsert().Model(c).Exec(context.Background())
	require.NoError(t, err)

	state := new(models.State)
	err = db.NewSelect().Model(state).Where("name = ?", "todo").Scan(context.Background())
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "test-slug",
		ContextID: &c.ID,
		Type:      "link",
		Value:     "https://example.com",
		StateID:   &state.ID,
	}
	_, err = db.NewInsert().Model(item).Exec(context.Background())
	require.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestMigration_ContextForeignKey(t *testing.T) {
	db := newTestBunDB(t)

	c := &models.Context{Name: "fk-test"}
	_, err := db.NewInsert().Model(c).Exec(context.Background())
	require.NoError(t, err)

	state := new(models.State)
	err = db.NewSelect().Model(state).Where("name = ?", "todo").Scan(context.Background())
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "fk-slug",
		ContextID: &c.ID,
		Type:      "pr",
		Value:     "https://github.com/a/b/pull/1",
		StateID:   &state.ID,
	}
	_, err = db.NewInsert().Model(item).Exec(context.Background())
	require.NoError(t, err)

	_, err = db.NewDelete().Model((*models.Context)(nil)).Where("id = ?", c.ID).Exec(context.Background())
	require.NoError(t, err)

	count, err := db.NewSelect().Model((*models.Item)(nil)).Where("slug = ?", "fk-slug").Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, count, "item should be cascade-deleted with context")
}

func TestMigration_Idempotent(t *testing.T) {
	db := newTestBunDB(t)

	sqldb := db.DB
	goose.SetDialect(store.GooseDialect)
	goose.SetLogger(goose.NopLogger())
	require.NoError(t, goose.Up(sqldb, "."))

	count, err := db.NewSelect().Model((*models.State)(nil)).Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 4, count, "migration should be idempotent")
}

func TestMigration_DownThenUp(t *testing.T) {
	db := newTestBunDB(t)

	sqldb := db.DB
	goose.SetDialect(store.GooseDialect)
	goose.SetLogger(goose.NopLogger())

	require.NoError(t, goose.Down(sqldb, "."))

	count, err := db.NewSelect().Model((*models.State)(nil)).Count(context.Background())
	assert.Error(t, err, "table should not exist after down migration")

	require.NoError(t, goose.Up(sqldb, "."))

	count, err = db.NewSelect().Model((*models.State)(nil)).Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 4, count, "states should be restored after re-running up")
	_ = count
}

func TestMigration_DefaultStatesNotOrphaned(t *testing.T) {
	db := newTestBunDB(t)

	var states []models.State
	err := db.NewSelect().Model(&states).Scan(context.Background())
	require.NoError(t, err)
	for _, s := range states {
		assert.False(t, s.Orphaned, "default state %q should not be orphaned", s.Name)
	}
}
