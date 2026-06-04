package cmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

// ── Rename command successes ─────────────────────────────────────────────────

func TestRenameCommand_Success(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "old", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "old"
	require.NoError(t, cfg.Save())

	buf := &bytes.Buffer{}
	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"old", "new"})
	require.NoError(t, renameCmd.Execute())

	assert.Contains(t, buf.String(), `"old" → "new"`)
	assert.Equal(t, "new", cfg.CurrentContext)
}

func TestRenameCommand_SuccessAliasRn(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "old", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"old", "new"})
	require.NoError(t, renameCmd.Execute())

	assert.Contains(t, buf.String(), `"old" → "new"`)
}

func TestRenameCommand_SuccessActiveContextUpdate(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "fix-auth", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "fix-auth"
	require.NoError(t, cfg.Save())

	buf := &bytes.Buffer{}
	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"fix-auth", "fix-auth-bug"})
	require.NoError(t, renameCmd.Execute())

	assert.Contains(t, buf.String(), `"fix-auth" → "fix-auth-bug"`)
	assert.Equal(t, "fix-auth-bug", cfg.CurrentContext)

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "fix-auth-bug", loaded.CurrentContext)
}

func TestRenameCommand_SuccessDoesNotAffectNonActive(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "background", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "other-context"
	require.NoError(t, cfg.Save())

	buf := &bytes.Buffer{}
	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"background", "renamed-bg"})
	require.NoError(t, renameCmd.Execute())

	assert.Contains(t, buf.String(), `"background" → "renamed-bg"`)
	assert.Equal(t, "other-context", cfg.CurrentContext)
}

// ── Rename command errors ────────────────────────────────────────────────────

func TestRenameCommand_ErrorNewNameExists(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "a", "")
	s.CreateContext(ctx, "b", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"a", "b"})
	err := renameCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRenameCommand_ErrorEmptyNewName(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "old", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"old", ""})
	err := renameCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "new name cannot be empty")
}

func TestRenameCommand_ErrorSameName(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "same", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"same", "same"})
	err := renameCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "new name cannot be same as old name")
}

func TestRenameCommand_ErrorOldNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"nope", "yep"})
	err := renameCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database update matched no rows")
}

func TestRenameCommand_ErrorNoArgs(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{})
	err := renameCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 2 arg")
}

func TestRenameCommand_ErrorOneArg(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	renameCmd := cmd.NewRenameCmd(s, cfg, buf)
	renameCmd.SetArgs([]string{"only-one"})
	err := renameCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 2 arg")
}
