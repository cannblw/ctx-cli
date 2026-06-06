package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/store"
)

// NewRmCmd creates the `ctx rm` command.
func NewRmCmd(s *store.Store, cfg *config.Config, stdin io.Reader, stdout io.Writer) *cobra.Command {
	var (
		rmContextName string
		rmForce       bool
	)

	cmd := &cobra.Command{
		Use:     "rm",
		Aliases: []string{"remove", "delete"},
		Short:   "Delete a context or item",
		Long: `Delete a context and all its items, or delete a single item.

Examples:
  # Delete a context (with confirmation)
  ctx rm --context fix-auth

  # Delete a context (skip confirmation)
  ctx rm --context fix-auth --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if rmContextName != "" {
				return deleteContext(s, cfg, rmContextName, rmForce, stdin, stdout)
			}
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("item deletion not yet implemented — use --context to delete a context")
		},
	}

	cmd.Flags().StringVar(&rmContextName, "context", "", "name of the context to delete")
	cmd.Flags().BoolVarP(&rmForce, "force", "f", false, "skip confirmation prompt")
	cmd.SetOut(stdout)
	return cmd
}

func deleteContext(s *store.Store, cfg *config.Config, name string, force bool, stdin io.Reader, stdout io.Writer) error {
	_, err := s.GetContext(context.Background(), name)
	if err != nil {
		return fmt.Errorf("context %q not found", name)
	}

	if !force {
		fmt.Fprintf(stdout, "Delete context %q and all its items? [y/N] ", name)
		reader := bufio.NewReader(stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Fprintln(stdout, "Cancelled")
			return nil
		}
	}

	if err := s.DeleteContext(context.Background(), name); err != nil {
		return err
	}

	if cfg.CurrentContext == name {
		cfg.CurrentContext = ""
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("context deleted but could not update config: %w", err)
		}
		fmt.Fprintln(stdout, "(active context cleared)")
	}

	fmt.Fprintf(stdout, "Deleted context %q\n", name)
	return nil
}
