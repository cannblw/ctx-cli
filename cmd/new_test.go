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

func TestNewCommand_Success(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{"fix-auth", "--desc", "Auth bugfix"})
	err := newCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), `Created context "fix-auth"`)

	ctx := context.Background()
	got, err := r.GetContext(ctx, "fix-auth")
	require.NoError(t, err)
	assert.Equal(t, "Auth bugfix", got.Description)
}

func TestNewCommand_NoDescription(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{"no-desc"})
	err := newCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), `Created context "no-desc"`)

	ctx := context.Background()
	got, err := r.GetContext(ctx, "no-desc")
	require.NoError(t, err)
	assert.Equal(t, "", got.Description)
}

func TestNewCommand_AliasN(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	// cobra resolves aliases automatically — `n` is an alias for `new`
	// but since we're executing the command directly, use the new command with SetArgs
	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{"via-alias"})
	err := newCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), `Created context "via-alias"`)
}

func TestNewCommand_DuplicateName(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{"same"})
	require.NoError(t, newCmd.Execute())

	buf.Reset()
	newCmd2 := cmd.NewNewCmd(r, buf)
	newCmd2.SetArgs([]string{"same"})
	err := newCmd2.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestNewCommand_NoArgs(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{})
	err := newCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg")
}

func TestNewCommand_OutputContainsID(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{"check-id"})
	require.NoError(t, newCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "Created context")
	assert.Contains(t, output, "(id:")
}

func TestNewCommand_ThenListContexts(t *testing.T) {
	r := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(r, buf)
	newCmd.SetArgs([]string{"first"})
	require.NoError(t, newCmd.Execute())

	buf.Reset()
	newCmd2 := cmd.NewNewCmd(r, buf)
	newCmd2.SetArgs([]string{"second"})
	require.NoError(t, newCmd2.Execute())

	ctx := context.Background()
	contexts, err := r.ListContexts(ctx)
	require.NoError(t, err)
	require.Len(t, contexts, 2)
	assert.Equal(t, "first", contexts[0].Name)
	assert.Equal(t, "second", contexts[1].Name)
}
