package repo_test

import (
	"context"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/pkg/models"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

func findMigrationsDir() string {
	dirs := []string{"migrations", "../../migrations", "../../../migrations"}
	for _, d := range dirs {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			abspath, _ := os.Getwd()
			_ = abspath
			return d
		}
	}
	return "migrations"
}

func TestMigration_CreatesDefaultStates(t *testing.T) {
	db := testutil.NewTestDB(t)

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
	db := testutil.NewTestDB(t)

	sqldb := db.DB
	goose.SetDialect("sqlite3")
	require.NoError(t, goose.Up(sqldb, findMigrationsDir()))

	count, err := db.NewSelect().Model((*models.State)(nil)).Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 4, count, "migration should be idempotent")
}
