package cmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

// ── Contexts command successes ───────────────────────────────────────────────

func TestContextsCommand_SuccessEmpty(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{})
	require.NoError(t, ctxCmd.Execute())

	assert.Contains(t, buf.String(), "No contexts yet")
}

func TestContextsCommand_SuccessWithContexts(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "Fix the auth bug")
	require.NoError(t, err)
	_, err = s.CreateContext(ctx, "migrate-db", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{})
	require.NoError(t, ctxCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "fix-auth")
	assert.Contains(t, output, "migrate-db")
	assert.Contains(t, output, "Fix the auth bug")
}

func TestContextsCommand_SuccessAliasC(t *testing.T) {
	s := testutil.NewTestDB(t)

	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "Fix the auth bug")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{})
	require.NoError(t, ctxCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "fix-auth")
	assert.Contains(t, output, "Fix the auth bug")
}

// ── Contexts command switch-via-name ──────────────────────────────────────────

func TestContextsCommand_SuccessSwitchWithName(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{"fix-auth"})
	require.NoError(t, ctxCmd.Execute())

	assert.Contains(t, buf.String(), `Switched to context "fix-auth"`)
	assert.Equal(t, "fix-auth", cfg.ActiveContext)
}

func TestContextsCommand_ErrorSwitchEmptyName(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{""})
	err := ctxCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context name is required")
}

func TestContextsCommand_ErrorSwitchNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{"nonexistent"})
	err := ctxCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), `context "nonexistent" not found`)
}

// ── Contexts command errors ──────────────────────────────────────────────────

func TestContextsCommand_ErrorDBFailure(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.Close()

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	ctxCmd := cmd.NewContextsCmd(s, cfg, buf)
	ctxCmd.SetArgs([]string{})
	err := ctxCmd.Execute()

	assert.Error(t, err)
}
