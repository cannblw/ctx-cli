package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/store"
)

func NewSwitchCmd(s *store.Store, cfg *config.Config, stdout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch <name>",
		Aliases: []string{"s"},
		Short:   "Set the active context",
		Long: `Set the active context (workstream). Subsequent commands like 'ctx add'
and 'ctx ls' will operate on this context by default.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return switchContext(s, cfg, stdout, args[0])
		},
	}
	return cmd
}

func switchContext(s *store.Store, cfg *config.Config, stdout io.Writer, name string) error {
	if name == "" || name == "global" {
		if err := cfg.ClearCurrentContext(); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "Switched to global")
		return nil
	}

	ctx := context.Background()

	_, err := s.GetContext(ctx, name)
	if err != nil {
		return fmt.Errorf("context %q not found — create it with: ctx new %s", name, name)
	}

	if err := cfg.SetCurrentContext(name); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Switched to context %q\n", name)
	return nil
}
