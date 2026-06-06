package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/models"
	"github.com/cannblw/ctx-cli/pkg/slug"
	"github.com/cannblw/ctx-cli/pkg/store"
)

var validTypes = map[string]struct{}{
	"link":   {},
	"pr":     {},
	"ticket": {},
	"file":   {},
}

// NewAddCmd creates the `ctx add` command.
func NewAddCmd(s *store.Store, cfg *config.Config, stdout io.Writer) *cobra.Command {
	var (
		addType   string
		addGlobal bool
	)

	cmd := &cobra.Command{
		Use:     "add <value>",
		Aliases: []string{"a"},
		Short:   "Add an item to the current context",
		Long: `Add an item (link, PR, ticket, file path, etc.) to the current context.

Type is auto-detected from the value. Use --type to override.
Use --global to add the item globally.

Examples:
  ctx add https://github.com/org/repo/pull/42
  ctx add ~/notes/todo.md --type file
  ctx a https://github.com/org/repo/pull/99 --global`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := args[0]
			if value == "" {
				return fmt.Errorf("value cannot be empty")
			}

			itemType := addType
			if itemType == "" {
				itemType = DetectType(value)
			}

			if _, ok := validTypes[itemType]; !ok {
				return fmt.Errorf("invalid type %q: must be one of link, pr, ticket, file", itemType)
			}

			var contextID *int64
			if !addGlobal {
				currentCtx := cfg.CurrentContext
				if currentCtx == "" {
					currentCtx = config.GlobalContextName
				}
				c, err := s.GetContext(context.Background(), currentCtx)
				if err != nil {
					return fmt.Errorf("current context %q not found", currentCtx)
				}
				contextID = &c.ID
			}

			todoState, err := s.GetState(context.Background(), "todo")
			if err != nil {
				return fmt.Errorf("default 'todo' state not found: %w", err)
			}

			itemSlug := slug.FromValue(value)
			item := &models.Item{
				Slug:      itemSlug,
				ContextID: contextID,
				Type:      itemType,
				Value:     value,
				StateID:   &todoState.ID,
			}

			if err := createItemWithRetry(s, item); err != nil {
				return err
			}

			scope := "globally"
			if contextID != nil {
				scope = fmt.Sprintf("in context %q", cfg.CurrentContext)
			}
			fmt.Fprintf(stdout, "Added [%s] %s (%s) %s\n", item.Slug, value, itemType, scope)
			return nil
		},
	}

	cmd.Flags().StringVarP(&addType, "type", "t", "", "item type: link, pr, ticket, file")
	cmd.Flags().BoolVarP(&addGlobal, "global", "g", false, "add item globally (no context)")
	return cmd
}

// DetectType auto-detects the item type from the value.
func DetectType(value string) string {
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return "link"
	}
	if _, err := os.Stat(value); err == nil {
		return "file"
	}
	if strings.HasPrefix(value, "~/") {
		expanded, err := os.UserHomeDir()
		if err == nil {
			path := expanded + value[1:]
			if _, err := os.Stat(path); err == nil {
				return "file"
			}
		}
	}
	return "link"
}

func createItemWithRetry(s *store.Store, item *models.Item) error {
	baseSlug := item.Slug
	for i := 0; i < 10; i++ {
		err := s.CreateItem(context.Background(), item)
		if err == nil {
			return nil
		}
		if !strings.Contains(err.Error(), "UNIQUE constraint") {
			return err
		}
		item.Slug = fmt.Sprintf("%s-%d", baseSlug, i+2)
	}
	return fmt.Errorf("could not create item: slug collision after 10 retries for %q", baseSlug)
}
