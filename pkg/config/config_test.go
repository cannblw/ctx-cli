package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/pkg/config"
)

// ── Load ─────────────────────────────────────────────────────────────────────

func TestLoad_SuccessCreatesDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "", cfg.CurrentContext)
	assert.Equal(t, config.DefaultItemState, cfg.DefaultItemState)

	_, err = os.Stat(filepath.Join(tmp, ".ctx", "config.yaml"))
	require.NoError(t, err)
}

func TestLoad_SuccessReadsExisting(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfgDir := filepath.Join(tmp, ".ctx")
	require.NoError(t, os.MkdirAll(cfgDir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(cfgDir, "config.yaml"),
		[]byte("current_context: fix-auth\n"),
		0644,
	))

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "fix-auth", cfg.CurrentContext)
}

// ── Load errors ──────────────────────────────────────────────────────────────

func TestLoad_ErrorCorruptedYAML(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfgDir := filepath.Join(tmp, ".ctx")
	require.NoError(t, os.MkdirAll(cfgDir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(cfgDir, "config.yaml"),
		[]byte("{{{not yaml!!!\n"),
		0644,
	))

	_, err := config.Load()
	assert.Error(t, err)
}

// ── Save ─────────────────────────────────────────────────────────────────────

func TestSave_SuccessRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)

	cfg.CurrentContext = "migrate-db"
	require.NoError(t, cfg.Save())

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "migrate-db", loaded.CurrentContext)
}

func TestSave_SuccessClearsCurrentContext(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)

	cfg.CurrentContext = "temp-ctx"
	require.NoError(t, cfg.Save())

	cfg.CurrentContext = ""
	require.NoError(t, cfg.Save())

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Empty(t, loaded.CurrentContext)
}

// ── SetCurrentContext ────────────────────────────────────────────────────────

func TestSetCurrentContext_Success(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)
	require.NoError(t, cfg.SetCurrentContext("fix-auth"))

	assert.Equal(t, "fix-auth", cfg.CurrentContext)

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "fix-auth", loaded.CurrentContext)
}

// ── ClearCurrentContext ──────────────────────────────────────────────────────

func TestClearCurrentContext_Success(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)
	require.NoError(t, cfg.SetCurrentContext("fix-auth"))
	assert.Equal(t, "fix-auth", cfg.CurrentContext)

	require.NoError(t, cfg.ClearCurrentContext())
	assert.Equal(t, "", cfg.CurrentContext)

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Empty(t, loaded.CurrentContext)
}

// ── GetContextDir ────────────────────────────────────────────────────────────

func TestGetContextDir_Success(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := config.GetContextDir("my-ctx")
	assert.Contains(t, dir, ".ctx")
	assert.Contains(t, dir, "contexts")
	assert.Contains(t, dir, "my-ctx")
}

// ── DefaultState ─────────────────────────────────────────────────────────────

func TestDefaultState_SuccessRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)

	cfg.DefaultItemState = "done"
	require.NoError(t, cfg.Save())

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "done", loaded.DefaultItemState)
}
