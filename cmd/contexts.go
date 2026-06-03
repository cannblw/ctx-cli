package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/cmd/internal/format"
	"github.com/cannblw/ctx-cli/pkg/models"
	"github.com/cannblw/ctx-cli/pkg/store"
)

// NewContextsCmd creates the `ctx contexts` command.
func NewContextsCmd(s *store.Store, stdout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contexts",
		Aliases: []string{"c", "ctx"},
		Short:   "List all contexts",
		Long:    `List all tracked contexts (workstreams) with their ID, name, and description.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			contexts, err := s.ListContexts(ctx)
			if err != nil {
				return fmt.Errorf("could not list contexts: %w", err)
			}

			if len(contexts) == 0 {
				fmt.Fprintln(stdout, "No contexts yet. Create one with: ctx new <name>")
				return nil
			}

			format.PrintTable(stdout,
				[]string{"ID", "NAME", "DESCRIPTION"},
				contexts,
				func(c models.Context) []string {
					return []string{
						fmt.Sprintf("%d", c.ID),
						c.Name,
						format.EmptyToDash(c.Description),
					}
				})

			return nil
		},
	}
	return cmd
}
