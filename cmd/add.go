package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
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
		Use:     "add <item>",
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

			autoDetected := false
			itemType := addType
			if itemType == "" {
				var err error
				itemType, err = DetectType(value)
				if err != nil {
					return err
				}
				autoDetected = true
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

			itemSlug, err := nextSlug(s, slug.FromValue(value))
			if err != nil {
				return err
			}

			item := &models.Item{
				Slug:      itemSlug,
				ContextID: contextID,
				Type:      itemType,
				Value:     value,
				StateID:   &todoState.ID,
			}

			if err := s.CreateItem(context.Background(), item); err != nil {
				return err
			}

			typeLabel := itemType
			if autoDetected {
				typeLabel = "auto-detected as " + itemType
			}
			scope := "globally"
			if contextID != nil {
				scope = fmt.Sprintf("in context %q", cfg.CurrentContext)
			}
			fmt.Fprintf(stdout, "Added [%s] %s (%s) %s\n", item.Slug, value, typeLabel, scope)
			return nil
		},
	}

	cmd.Flags().StringVarP(&addType, "type", "t", "", "item type: link, pr, ticket, file")
	cmd.Flags().BoolVarP(&addGlobal, "global", "g", false, "add item globally (no context)")
	return cmd
}

// DetectType auto-detects the item type from the value.
func DetectType(value string) (string, error) {
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return "link", nil
	}
	if _, err := os.Stat(value); err == nil {
		return "file", nil
	}
	if strings.HasPrefix(value, "~/") {
		expanded, err := os.UserHomeDir()
		if err == nil {
			path := expanded + value[1:]
			if _, err := os.Stat(path); err == nil {
				return "file", nil
			}
		}
	}
	return "", fmt.Errorf("could not detect type for %q, use --type to specify one of link, pr, ticket, file", value)
}

func nextSlug(s *store.Store, base string) (string, error) {
	existing, err := s.FindSlugsByPrefix(context.Background(), base)
	if err != nil {
		return "", err
	}
	if len(existing) == 0 {
		return base, nil
	}

	maxSuffix := 1
	for _, slug := range existing {
		if slug == base {
			continue
		}
		rest := strings.TrimPrefix(slug, base+"-")
		if n, err := strconv.Atoi(rest); err == nil && n > maxSuffix {
			maxSuffix = n
		}
	}
	return fmt.Sprintf("%s-%d", base, maxSuffix+1), nil
}
