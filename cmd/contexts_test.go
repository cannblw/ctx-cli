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
	buf := &bytes.Buffer{}

	ctxCmd := cmd.NewContextsCmd(s, buf)
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

	buf := &bytes.Buffer{}
	ctxCmd := cmd.NewContextsCmd(s, buf)
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

	buf := &bytes.Buffer{}
	ctxCmd := cmd.NewContextsCmd(s, buf)
	ctxCmd.SetArgs([]string{})
	require.NoError(t, ctxCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "fix-auth")
	assert.Contains(t, output, "Fix the auth bug")
}

// ── Contexts command errors ──────────────────────────────────────────────────

func TestContextsCommand_ErrorDBFailure(t *testing.T) {
	s := testutil.NewTestDB(t)
	s.Close()

	buf := &bytes.Buffer{}
	ctxCmd := cmd.NewContextsCmd(s, buf)
	ctxCmd.SetArgs([]string{})
	err := ctxCmd.Execute()

	assert.Error(t, err)
}
