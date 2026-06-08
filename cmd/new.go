package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/pkg/store"
)

func NewNewCmd(s *store.Store, stdout io.Writer) *cobra.Command {
	var descFlag string

	cmd := &cobra.Command{
		Use:     "new <name>",
		Aliases: []string{"n"},
		Short:   "Create a new context",
		Long: `Create a new context (workstream) to track items under.

Examples:
  ctx new fix-auth-bug
  ctx n upgrade-auth --description "Refactor OIDC flow to support multi-tenant"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if name == "" || name == "global" {
				return fmt.Errorf("context name cannot be empty or \"global\"")
			}

			ctx := context.Background()
			c, err := s.CreateContext(ctx, name, descFlag)
			if err != nil {
				return fmt.Errorf("could not create context: %w", err)
			}

			fmt.Fprintf(stdout, "Created context %q (id: %d)\n", c.Name, c.ID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&descFlag, "description", "d", "", "description for the context")
	return cmd
}
