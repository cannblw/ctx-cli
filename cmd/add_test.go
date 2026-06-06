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
	"github.com/cannblw/ctx-cli/pkg/slug"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

// ── Add command successes ────────────────────────────────────────────────────

func TestAddCommand_Success(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"https://github.com/org/repo/pull/42"})
	require.NoError(t, addCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "Added")
	assert.Contains(t, output, "https://github.com/org/repo/pull/42")
	assert.Contains(t, output, "auto-detected as link")
	assert.Contains(t, output, `in context "my-ctx"`)
}

func TestAddCommand_SuccessGlobal(t *testing.T) {
	s := testutil.NewTestDB(t)

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"https://example.com", "--global"})
	require.NoError(t, addCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "Added")
	assert.Contains(t, output, "auto-detected as link")
	assert.Contains(t, output, "globally")
}

func TestAddCommand_SuccessAliasA(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"https://example.com"})
	require.NoError(t, addCmd.Execute())

	assert.Contains(t, buf.String(), "auto-detected as link")
}

func TestAddCommand_SuccessTypeOverride(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"https://jira.example.com/PROJ-123", "--type", "ticket"})
	require.NoError(t, addCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "(ticket)")
	assert.NotContains(t, output, "auto-detected")
}

func TestAddCommand_SuccessFileType(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	tmp := t.TempDir()
	f := filepath.Join(tmp, "notes.md")
	require.NoError(t, os.WriteFile(f, []byte("test"), 0644))

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{f})
	require.NoError(t, addCmd.Execute())

	assert.Contains(t, buf.String(), "auto-detected as file")
}

func TestAddCommand_SuccessFileCopyToContextDir(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	tmp := t.TempDir()
	f := filepath.Join(tmp, "notes.md")
	require.NoError(t, os.WriteFile(f, []byte("hello"), 0644))

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{f})
	require.NoError(t, addCmd.Execute())

	itemSlug := slug.FromValue(f)
	item, err := s.GetItem(ctx, itemSlug)
	require.NoError(t, err)
	assert.Contains(t, item.Value, ".ctx")
	assert.Contains(t, item.Value, "contexts")
	assert.Contains(t, item.Value, "my-ctx")
	assert.Contains(t, item.Value, "notes.md")

	copied, err := os.Stat(item.Value)
	require.NoError(t, err)
	assert.False(t, copied.IsDir())
}

func TestAddCommand_SuccessFileCopyGlobal(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()

	tmp := t.TempDir()
	f := filepath.Join(tmp, "notes.md")
	require.NoError(t, os.WriteFile(f, []byte("hello"), 0644))

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{f, "--global"})
	require.NoError(t, addCmd.Execute())

	itemSlug := slug.FromValue(f)
	item, err := s.GetItem(ctx, itemSlug)
	require.NoError(t, err)
	assert.Nil(t, item.ContextID)
	assert.Contains(t, item.Value, "contexts")
	assert.Contains(t, item.Value, "global")
	assert.Contains(t, item.Value, "notes.md")
}

func TestAddCommand_SuccessExplicitFileType(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	tmp := t.TempDir()
	f := filepath.Join(tmp, "notes.md")
	require.NoError(t, os.WriteFile(f, []byte("hello"), 0644))

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{f, "--type", "file"})
	require.NoError(t, addCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "(file)")
	assert.NotContains(t, output, "auto-detected")

	itemSlug := slug.FromValue(f)
	item, err := s.GetItem(ctx, itemSlug)
	require.NoError(t, err)
	assert.Contains(t, item.Value, ".ctx")
	assert.Contains(t, item.Value, "contexts")
}

func TestAddCommand_SuccessSlugCollision(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	_, err := s.CreateContext(ctx, "my-ctx", "")
	require.NoError(t, err)

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"https://example.com"})
	require.NoError(t, addCmd.Execute())

	item, err := s.GetItem(ctx, "https-example-com")
	require.NoError(t, err)
	assert.NotZero(t, item.ID)

	buf.Reset()
	addCmd2 := cmd.NewAddCmd(s, cfg, buf)
	addCmd2.SetArgs([]string{"https://example.com"})
	require.NoError(t, addCmd2.Execute())

	output := buf.String()
	assert.Contains(t, output, "Added")
	assert.Contains(t, output, "https-example-com-2")
}

// ── Add command errors ───────────────────────────────────────────────────────

func TestAddCommand_ErrorUndetectableType(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"just-some-text"})
	err := addCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not detect type")
	assert.Contains(t, err.Error(), "--type")
}

func TestAddCommand_ErrorInvalidType(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"foo", "--type", "invalid"})
	err := addCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid type")
}

func TestAddCommand_ErrorEmptyValue(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{""})
	err := addCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value cannot be empty")
}

func TestAddCommand_ErrorNoArgs(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{})
	err := addCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg")
}

func TestAddCommand_ErrorFileNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	ctx := context.Background()
	s.CreateContext(ctx, "my-ctx", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "my-ctx"
	buf := &bytes.Buffer{}

	addCmd := cmd.NewAddCmd(s, cfg, buf)
	addCmd.SetArgs([]string{"/nonexistent/file.txt", "--type", "file"})
	err := addCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not open file")
}

// ── DetectType ───────────────────────────────────────────────────────────────

func TestDetectType_SuccessURL(t *testing.T) {
	typ, err := cmd.DetectType("https://example.com")
	require.NoError(t, err)
	assert.Equal(t, "link", typ)

	typ, err = cmd.DetectType("http://example.com")
	require.NoError(t, err)
	assert.Equal(t, "link", typ)
}

func TestDetectType_SuccessExistingFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "notes.md")
	os.WriteFile(f, []byte("test"), 0644)

	typ, err := cmd.DetectType(f)
	require.NoError(t, err)
	assert.Equal(t, "file", typ)
}

func TestDetectType_ErrorPlainString(t *testing.T) {
	_, err := cmd.DetectType("just-some-text")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not detect type")
}
