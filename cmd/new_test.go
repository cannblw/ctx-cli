package cmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/repo"
	"github.com/cannblw/ctx-cli/pkg/testutil"
)

func TestNewCommand_Success(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)
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

func TestNewCommand_DuplicateName(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repo.NewBunRepositoryFromDB(db)
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
