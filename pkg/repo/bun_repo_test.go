package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/pkg/repo"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

func TestCreateContext(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)

	ctx := context.Background()
	c, err := r.CreateContext(ctx, "fix-auth", "Fix the auth bug")
	require.NoError(t, err)
	assert.Equal(t, "fix-auth", c.Name)
	assert.Equal(t, "Fix the auth bug", c.Description)
	assert.NotZero(t, c.ID)
	assert.False(t, c.CreatedAt.IsZero())
}

func TestCreateContext_DuplicateName(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)

	ctx := context.Background()
	_, err := r.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	_, err = r.CreateContext(ctx, "fix-auth", "")
	assert.Error(t, err)
}

func TestGetContext(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)

	ctx := context.Background()
	created, _ := r.CreateContext(ctx, "fix-auth", "desc")

	got, err := r.GetContext(ctx, "fix-auth")
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "fix-auth", got.Name)
}

func TestGetContext_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)

	_, err := r.GetContext(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestListContexts_Empty(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)

	contexts, err := r.ListContexts(context.Background())
	require.NoError(t, err)
	assert.Empty(t, contexts)
}

func TestListContexts_Ordered(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)

	ctx := context.Background()
	r.CreateContext(ctx, "b", "")
	r.CreateContext(ctx, "a", "")
	r.CreateContext(ctx, "c", "")

	contexts, err := r.ListContexts(ctx)
	require.NoError(t, err)
	require.Len(t, contexts, 3)
	assert.Equal(t, "b", contexts[0].Name)
	assert.Equal(t, "a", contexts[1].Name)
	assert.Equal(t, "c", contexts[2].Name)
}
