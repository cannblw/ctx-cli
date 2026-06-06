package cmd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

// ── resolveContextID ─────────────────────────────────────────────────────────

func TestResolveContextID_SuccessReturnsNilForGlobal(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)

	id, err := cmd.ResolveContextID(s, cfg, true)
	require.NoError(t, err)
	assert.Nil(t, id)
}

func TestResolveContextID_SuccessReturnsID(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	c, err := s.CreateContext(ctx, "my-ctx", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"

	id, err := cmd.ResolveContextID(s, cfg, false)
	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, c.ID, *id)
}

func TestResolveContextID_ErrorContextNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	cfg.CurrentContext = "nonexistent"

	_, err := cmd.ResolveContextID(s, cfg, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
