package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/pkg/models"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

// ── Context creation ─────────────────────────────────────────────────────────

func TestCreateContext_Success(t *testing.T) {
	s := testutil.NewTestDB(t)

	c, err := s.CreateContext(context.Background(), "fix-auth", "Fix the auth bug")
	require.NoError(t, err)
	assert.Equal(t, "fix-auth", c.Name)
	assert.Equal(t, "Fix the auth bug", c.Description)
	assert.NotZero(t, c.ID)
	assert.False(t, c.CreatedAt.IsZero())
	assert.False(t, c.UpdatedAt.IsZero())
	assert.True(t, c.UpdatedAt.Equal(c.CreatedAt), "created_at and updated_at should be equal on creation")
}

func TestCreateContext_SuccessEmptyDescription(t *testing.T) {
	s := testutil.NewTestDB(t)

	c, err := s.CreateContext(context.Background(), "empty-desc", "")
	require.NoError(t, err)
	assert.Equal(t, "empty-desc", c.Name)
	assert.Equal(t, "", c.Description)
	assert.NotZero(t, c.ID)
}

// ── Context retrieval ────────────────────────────────────────────────────────

func TestGetContext_Success(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	created, _ := s.CreateContext(ctx, "fix-auth", "desc")

	got, err := s.GetContext(ctx, "fix-auth")
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "fix-auth", got.Name)
	assert.Equal(t, "desc", got.Description)
}

// ── Context listing ──────────────────────────────────────────────────────────

func TestListContexts_SuccessEmpty(t *testing.T) {
	s := testutil.NewTestDB(t)

	contexts, err := s.ListContexts(context.Background())
	require.NoError(t, err)
	assert.Empty(t, contexts)
}

func TestListContexts_SuccessOneItem(t *testing.T) {
	s := testutil.NewTestDB(t)

	c, err := s.CreateContext(context.Background(), "solo", "only one")
	require.NoError(t, err)

	contexts, err := s.ListContexts(context.Background())
	require.NoError(t, err)
	require.Len(t, contexts, 1)
	assert.Equal(t, c.ID, contexts[0].ID)
	assert.Equal(t, "solo", contexts[0].Name)
	assert.Equal(t, "only one", contexts[0].Description)
}

func TestListContexts_SuccessOrdered(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	s.CreateContext(ctx, "b", "")
	s.CreateContext(ctx, "a", "")
	s.CreateContext(ctx, "c", "")

	contexts, err := s.ListContexts(ctx)
	require.NoError(t, err)
	require.Len(t, contexts, 3)
	assert.Equal(t, "b", contexts[0].Name)
	assert.Equal(t, "a", contexts[1].Name)
	assert.Equal(t, "c", contexts[2].Name)
}

// ── Context deletion ─────────────────────────────────────────────────────────

func TestContext_SuccessDeleteCascadesItems(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	c, err := s.CreateContext(ctx, "fk-test", "")
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "cascade-slug",
		ContextID: &c.ID,
		Type:      "link",
		Value:     "https://example.com",
		StateID:   nil,
	}
	require.NoError(t, s.CreateItem(ctx, item))

	err = s.DeleteContext(ctx, c.Name)
	require.NoError(t, err)

	contexts, err := s.ListContexts(ctx)
	require.NoError(t, err)
	assert.Empty(t, contexts, "context should be deleted")

	_, err = s.GetItem(ctx, "cascade-slug")
	assert.Error(t, err, "item should be cascade-deleted")
}

// ── Context creation errors ──────────────────────────────────────────────────

func TestCreateContext_ErrorDuplicateName(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	_, err = s.CreateContext(ctx, "fix-auth", "")
	assert.Error(t, err)
}

// ── Context retrieval errors ─────────────────────────────────────────────────

func TestGetContext_ErrorNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)

	_, err := s.GetContext(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestGetContext_ErrorCaseSensitive(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	s.CreateContext(ctx, "MyCtx", "")

	_, err := s.GetContext(ctx, "myctx")
	assert.Error(t, err, "lookup should be case-sensitive")

	got, err := s.GetContext(ctx, "MyCtx")
	require.NoError(t, err)
	assert.Equal(t, "MyCtx", got.Name)
}
