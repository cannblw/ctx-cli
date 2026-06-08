package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/store"
)

// NewRenameCmd creates the `ctx rename` command.
func NewRenameCmd(s *store.Store, cfg *config.Config, stdout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rename <old-name> <new-name>",
		Aliases: []string{"rn"},
		Short:   "Rename a context",
		Long: `Rename a context. If the active context is renamed, the active context
is updated automatically.

Examples:
  ctx rename fix-auth fix-auth-bug
  ctx rn old-name new-name`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName, newName := args[0], args[1]
			if oldName == "" || oldName == "global" {
				return fmt.Errorf("cannot rename the global context")
			}
			if newName == "" {
				return fmt.Errorf("new name cannot be empty")
			}
			if oldName == newName {
				return fmt.Errorf("new name cannot be same as old name")
			}

			ctx := context.Background()
			c, err := s.RenameContext(ctx, oldName, newName)
			if err != nil {
				return err
			}

			if cfg.CurrentContext == oldName {
				if err := cfg.SetCurrentContext(newName); err != nil {
					return fmt.Errorf("rename succeeded but could not update current context: %w", err)
				}
			}

			fmt.Fprintf(stdout, "Renamed context %q → %q (id: %d)\n", oldName, c.Name, c.ID)
			return nil
		},
	}
	return cmd
}
