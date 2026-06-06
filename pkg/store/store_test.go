package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/pkg/config"
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
	require.Len(t, contexts, 1)
	assert.Equal(t, config.GlobalContextName, contexts[0].Name)
}

func TestListContexts_SuccessOneItem(t *testing.T) {
	s := testutil.NewTestDB(t)

	c, err := s.CreateContext(context.Background(), "solo", "only one")
	require.NoError(t, err)

	contexts, err := s.ListContexts(context.Background())
	require.NoError(t, err)
	require.Len(t, contexts, 2)
	assert.Equal(t, config.GlobalContextName, contexts[0].Name)
	assert.Equal(t, c.ID, contexts[1].ID)
	assert.Equal(t, "solo", contexts[1].Name)
	assert.Equal(t, "only one", contexts[1].Description)
}

func TestListContexts_SuccessOrdered(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	s.CreateContext(ctx, "b", "")
	s.CreateContext(ctx, "a", "")
	s.CreateContext(ctx, "c", "")

	contexts, err := s.ListContexts(ctx)
	require.NoError(t, err)
	require.Len(t, contexts, 4)
	assert.Equal(t, config.GlobalContextName, contexts[0].Name)
	assert.Equal(t, "b", contexts[1].Name)
	assert.Equal(t, "a", contexts[2].Name)
	assert.Equal(t, "c", contexts[3].Name)
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
	assert.Len(t, contexts, 1, "only global context should remain")
	assert.Equal(t, config.GlobalContextName, contexts[0].Name)

	_, err = s.GetItem(ctx, "cascade-slug")
	assert.Error(t, err, "item should be cascade-deleted")
}

func TestDeleteContext_Success(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	s.CreateContext(ctx, "temp", "")
	require.NoError(t, s.DeleteContext(ctx, "temp"))

	_, err := s.GetContext(ctx, "temp")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get context")
}

func TestDeleteContext_ErrorNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)

	err := s.DeleteContext(context.Background(), "nope")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find")
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

// ── Context renaming ─────────────────────────────────────────────────────────

func TestRenameContext_Success(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	s.CreateContext(ctx, "old-name", "desc")

	c, err := s.RenameContext(ctx, "old-name", "new-name")
	require.NoError(t, err)
	assert.Equal(t, "new-name", c.Name)
	assert.Equal(t, "desc", c.Description)

	_, err = s.GetContext(ctx, "old-name")
	assert.Error(t, err)

	got, err := s.GetContext(ctx, "new-name")
	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
}

// ── Context renaming errors ──────────────────────────────────────────────────

func TestRenameContext_ErrorNewNameExists(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	s.CreateContext(ctx, "a", "")
	s.CreateContext(ctx, "b", "")

	_, err := s.RenameContext(ctx, "a", "b")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRenameContext_ErrorOldNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)

	_, err := s.RenameContext(context.Background(), "nope", "yep")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database update matched no rows")
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

// ── Item creation ────────────────────────────────────────────────────────────

func TestCreateItem_Success(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	c, err := s.CreateContext(ctx, "my-ctx", "")
	require.NoError(t, err)

	todo, err := s.GetState(ctx, config.DefaultItemState)
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "test-item",
		ContextID: &c.ID,
		Type:      "link",
		Value:     "https://example.com",
		StateID:   &todo.ID,
	}

	require.NoError(t, s.CreateItem(ctx, item))
	assert.NotZero(t, item.ID)
	assert.False(t, item.CreatedAt.IsZero())
}

func TestCreateItem_SuccessGlobal(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	todo, err := s.GetState(ctx, config.DefaultItemState)
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "global-item",
		ContextID: nil,
		Type:      "link",
		Value:     "https://example.com",
		StateID:   &todo.ID,
	}

	require.NoError(t, s.CreateItem(ctx, item))
	assert.NotZero(t, item.ID)
	assert.Nil(t, item.ContextID)
}

func TestCreateItem_ErrorDuplicateSlug(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	c, err := s.CreateContext(ctx, "my-ctx", "")
	require.NoError(t, err)

	todo, err := s.GetState(ctx, config.DefaultItemState)
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "dup-slug",
		ContextID: &c.ID,
		Type:      "link",
		Value:     "https://example.com",
		StateID:   &todo.ID,
	}
	require.NoError(t, s.CreateItem(ctx, item))

	dup := &models.Item{
		Slug:      "dup-slug",
		ContextID: &c.ID,
		Type:      "link",
		Value:     "https://other.com",
		StateID:   &todo.ID,
	}
	err = s.CreateItem(ctx, dup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not insert item")
}

// ── State retrieval ──────────────────────────────────────────────────────────

func TestGetState_SuccessTodo(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	st, err := s.GetState(ctx, config.DefaultItemState)
	require.NoError(t, err)
	assert.Equal(t, config.DefaultItemState, st.Name)
	assert.Equal(t, 0, st.Position)
}

func TestGetState_SuccessInProgress(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	st, err := s.GetState(ctx, "in-progress")
	require.NoError(t, err)
	assert.Equal(t, "in-progress", st.Name)
	assert.Equal(t, 1, st.Position)
}

func TestGetState_ErrorNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	_, err := s.GetState(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get state")
}

// ── Item retrieval ───────────────────────────────────────────────────────────

func TestGetItem_Success(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	c, err := s.CreateContext(ctx, "my-ctx", "")
	require.NoError(t, err)

	todo, err := s.GetState(ctx, config.DefaultItemState)
	require.NoError(t, err)

	item := &models.Item{
		Slug:      "get-item-test",
		ContextID: &c.ID,
		Type:      "file",
		Value:     "~/notes/foo.md",
		StateID:   &todo.ID,
	}
	require.NoError(t, s.CreateItem(ctx, item))

	got, err := s.GetItem(ctx, "get-item-test")
	require.NoError(t, err)
	assert.Equal(t, item.ID, got.ID)
	assert.Equal(t, "get-item-test", got.Slug)
	assert.Equal(t, "file", got.Type)
	assert.Equal(t, "~/notes/foo.md", got.Value)
	assert.Equal(t, c.ID, *got.ContextID)
}

func TestGetItem_ErrorNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	_, err := s.GetItem(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get item")
}
