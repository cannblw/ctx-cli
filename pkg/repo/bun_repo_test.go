package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/pkg/testutil"
)

func TestCreateContext(t *testing.T) {
	r := testutil.NewTestDB(t)

	ctx := context.Background()
	c, err := r.CreateContext(ctx, "fix-auth", "Fix the auth bug")
	require.NoError(t, err)
	assert.Equal(t, "fix-auth", c.Name)
	assert.Equal(t, "Fix the auth bug", c.Description)
	assert.NotZero(t, c.ID)
	assert.False(t, c.CreatedAt.IsZero())
	assert.False(t, c.UpdatedAt.IsZero())
}

func TestCreateContext_EmptyDescription(t *testing.T) {
	r := testutil.NewTestDB(t)

	c, err := r.CreateContext(context.Background(), "empty-desc", "")
	require.NoError(t, err)
	assert.Equal(t, "empty-desc", c.Name)
	assert.Equal(t, "", c.Description)
	assert.NotZero(t, c.ID)
}

func TestCreateContext_DuplicateName(t *testing.T) {
	r := testutil.NewTestDB(t)

	ctx := context.Background()
	_, err := r.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	_, err = r.CreateContext(ctx, "fix-auth", "")
	assert.Error(t, err)
}

func TestGetContext(t *testing.T) {
	r := testutil.NewTestDB(t)

	ctx := context.Background()
	created, _ := r.CreateContext(ctx, "fix-auth", "desc")

	got, err := r.GetContext(ctx, "fix-auth")
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "fix-auth", got.Name)
	assert.Equal(t, "desc", got.Description)
}

func TestGetContext_NotFound(t *testing.T) {
	r := testutil.NewTestDB(t)

	_, err := r.GetContext(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestGetContext_AfterCreate(t *testing.T) {
	r := testutil.NewTestDB(t)

	ctx := context.Background()
	c, err := r.CreateContext(ctx, "fresh", "just created")
	require.NoError(t, err)

	got, err := r.GetContext(ctx, "fresh")
	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	assert.Equal(t, c.CreatedAt, got.CreatedAt)
}

func TestListContexts_Empty(t *testing.T) {
	r := testutil.NewTestDB(t)

	contexts, err := r.ListContexts(context.Background())
	require.NoError(t, err)
	assert.Empty(t, contexts)
}

func TestListContexts_SingleItem(t *testing.T) {
	r := testutil.NewTestDB(t)

	c, err := r.CreateContext(context.Background(), "solo", "only one")
	require.NoError(t, err)

	contexts, err := r.ListContexts(context.Background())
	require.NoError(t, err)
	require.Len(t, contexts, 1)
	assert.Equal(t, c.ID, contexts[0].ID)
	assert.Equal(t, "solo", contexts[0].Name)
	assert.Equal(t, "only one", contexts[0].Description)
}

func TestListContexts_Ordered(t *testing.T) {
	r := testutil.NewTestDB(t)

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

func TestCreateContext_UpdatedAtEqualsCreatedAtInitially(t *testing.T) {
	r := testutil.NewTestDB(t)

	c, err := r.CreateContext(context.Background(), "time-check", "")
	require.NoError(t, err)
	assert.True(t, c.UpdatedAt.Equal(c.CreatedAt), "created_at and updated_at should be equal on creation")
}

func TestClose(t *testing.T) {
	r := testutil.NewTestDB(t)

	assert.NoError(t, r.Close())
}

func TestGetContext_CaseSensitive(t *testing.T) {
	r := testutil.NewTestDB(t)

	ctx := context.Background()
	r.CreateContext(ctx, "MyCtx", "")

	_, err := r.GetContext(ctx, "myctx")
	assert.Error(t, err, "lookup should be case-sensitive")

	got, err := r.GetContext(ctx, "MyCtx")
	require.NoError(t, err)
	assert.Equal(t, "MyCtx", got.Name)
}
