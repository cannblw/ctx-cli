package cmd_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

func setupConfig(t *testing.T) *config.Config {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	cfg, err := config.Load()
	require.NoError(t, err)
	return cfg
}

func readConfigFile(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(config.Dir(), "config.yaml"))
	require.NoError(t, err)
	return string(data)
}

// ── Switch command successes ──────────────────────────────────────────────────

func TestSwitchCommand_SuccessSwitchToContext(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"fix-auth"})
	require.NoError(t, switchCmd.Execute())

	assert.Contains(t, buf.String(), `Switched to context "fix-auth"`)
	assert.Equal(t, "fix-auth", cfg.CurrentContext)

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "fix-auth", loaded.CurrentContext)
}

func TestSwitchCommand_ErrorEmptyName(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{""})
	err := switchCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context name is required")
}

func TestSwitchCommand_SuccessAliasS(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "my-ctx", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"my-ctx"})
	require.NoError(t, switchCmd.Execute())

	assert.Contains(t, buf.String(), `Switched to context "my-ctx"`)
}

func TestSwitchCommand_SuccessReswitchSameContext(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"fix-auth"})
	require.NoError(t, switchCmd.Execute())

	buf.Reset()
	switchCmd = cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"fix-auth"})
	require.NoError(t, switchCmd.Execute())

	assert.Contains(t, buf.String(), `Switched to context "fix-auth"`)
}

func TestSwitchCommand_SuccessSwitchToGlobal(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{config.GlobalContextName})
	require.NoError(t, switchCmd.Execute())

	assert.Contains(t, buf.String(), `Switched to context "`+config.GlobalContextName+`"`)
	assert.Equal(t, config.GlobalContextName, cfg.CurrentContext)

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, config.GlobalContextName, loaded.CurrentContext)
}

// ── Switch command errors ─────────────────────────────────────────────────────

func TestSwitchCommand_ErrorSaveConfigFails(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	require.NoError(t, os.RemoveAll(config.Dir()))

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"fix-auth"})
	err = switchCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not save config")
}

func TestSwitchCommand_ErrorNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"nonexistent"})
	err := switchCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), `context "nonexistent" not found`)
}

func TestSwitchCommand_ErrorNoArgs(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{})
	err := switchCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg")
}

// ── Config persistence ────────────────────────────────────────────────────────

func TestSwitchCommand_SuccessPersistedToFile(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "fix-auth", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	switchCmd := cmd.NewSwitchCmd(s, cfg, buf)
	switchCmd.SetArgs([]string{"fix-auth"})
	require.NoError(t, switchCmd.Execute())

	data := readConfigFile(t)
	assert.Contains(t, data, "current_context: fix-auth")
}
