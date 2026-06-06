package cmd_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

// ── rm command successes ─────────────────────────────────────────────────────

func TestRmCommand_SuccessForce(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "temp"
	require.NoError(t, cfg.Save())

	buf := &bytes.Buffer{}
	stdin := strings.NewReader("")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp", "--force"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), `Deleted context "temp"`)
	assert.Contains(t, buf.String(), "(active context cleared)")
	assert.Equal(t, "", cfg.CurrentContext)

	_, err := s.GetContext(context.Background(), "temp")
	assert.Error(t, err)
}

func TestRmCommand_SuccessForceNotActive(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	cfg.CurrentContext = "other"
	require.NoError(t, cfg.Save())

	buf := &bytes.Buffer{}
	stdin := strings.NewReader("")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp", "--force"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), `Deleted context "temp"`)
	assert.NotContains(t, buf.String(), "(active context cleared)")
	assert.Equal(t, "other", cfg.CurrentContext)
}

func TestRmCommand_SuccessConfirmYes(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	stdin := strings.NewReader("y\n")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), `Deleted context "temp"`)
}

func TestRmCommand_SuccessAliasRemove(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	stdin := strings.NewReader("")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp", "--force"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), `Deleted context "temp"`)
}

func TestRmCommand_SuccessAliasDelete(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	stdin := strings.NewReader("")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp", "--force"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), `Deleted context "temp"`)
}

func TestRmCommand_SuccessCancelledNo(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	stdin := strings.NewReader("n\n")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), "Cancelled")

	_, err := s.GetContext(context.Background(), "temp")
	assert.NoError(t, err)
}

func TestRmCommand_SuccessCancelledEmpty(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.CreateContext(context.Background(), "temp", "")

	cfg := setupConfig(t)
	buf := &bytes.Buffer{}
	stdin := strings.NewReader("\n")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "temp"})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), "Cancelled")

	_, err := s.GetContext(context.Background(), "temp")
	assert.NoError(t, err)
}

func TestRmCommand_SuccessHelpNoArgs(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)

	buf := &bytes.Buffer{}
	stdin := strings.NewReader("")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{})
	require.NoError(t, rmCmd.Execute())

	assert.Contains(t, buf.String(), "Delete a context and all its items")
}

// ── rm command errors ────────────────────────────────────────────────────────

func TestRmCommand_ErrorNotFound(t *testing.T) {
	s := testutil.NewTestDB(t)
	cfg := setupConfig(t)

	buf := &bytes.Buffer{}
	stdin := strings.NewReader("")

	rmCmd := cmd.NewRmCmd(s, cfg, stdin, buf)
	rmCmd.SetArgs([]string{"--context", "nope", "--force"})
	err := rmCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), `could not find context "nope"`)
}
