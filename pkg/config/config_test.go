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
	assert.Equal(t, config.GlobalContextName, cfg.CurrentContext)

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
