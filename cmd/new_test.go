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
	s := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(s, buf)
	newCmd.SetArgs([]string{"fix-auth"})
	err := newCmd.Execute()

	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, `Created context "fix-auth"`)
	assert.Contains(t, output, "(id:")

	ctx := context.Background()
	got, err := s.GetContext(ctx, "fix-auth")
	require.NoError(t, err)
	assert.Equal(t, "", got.Description)
}

func TestNewCommand_SuccessWithDescription(t *testing.T) {
	s := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(s, buf)
	newCmd.SetArgs([]string{"fix-auth", "--description", "Auth bugfix"})
	err := newCmd.Execute()

	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, `Created context "fix-auth"`)
	assert.Contains(t, output, "(id:")

	ctx := context.Background()
	got, err := s.GetContext(ctx, "fix-auth")
	require.NoError(t, err)
	assert.Equal(t, "Auth bugfix", got.Description)
}

func TestNewCommand_SuccessAliasN(t *testing.T) {
	s := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(s, buf)
	newCmd.SetArgs([]string{"via-alias"})
	err := newCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), `Created context "via-alias"`)

	ctx := context.Background()
	got, err := s.GetContext(ctx, "via-alias")
	require.NoError(t, err)
	assert.Equal(t, "via-alias", got.Name)
}

func TestNewCommand_ErrorNoArgs(t *testing.T) {
	s := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(s, buf)
	newCmd.SetArgs([]string{})
	err := newCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg")
}

func TestNewCommand_ErrorDuplicateName(t *testing.T) {
	s := testutil.NewTestDB(t)
	buf := new(bytes.Buffer)

	newCmd := cmd.NewNewCmd(s, buf)
	newCmd.SetArgs([]string{"same"})
	require.NoError(t, newCmd.Execute())

	buf.Reset()
	newCmd2 := cmd.NewNewCmd(s, buf)
	newCmd2.SetArgs([]string{"same"})
	err := newCmd2.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}
