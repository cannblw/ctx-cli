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

			itemType, autoDetected, err := resolveItemType(value, addType)
			if err != nil {
				return err
			}

			contextID, err := resolveContextID(s, cfg, addGlobal)
			if err != nil {
				return err
			}

			item, err := insertItem(s, value, itemType, contextID)
			if err != nil {
				return err
			}

			printAdded(stdout, item, value, itemType, autoDetected, contextID, cfg)
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

func resolveItemType(value, addType string) (string, bool, error) {
	if addType != "" {
		if _, ok := validTypes[addType]; !ok {
			return "", false, fmt.Errorf("invalid type %q: must be one of link, pr, ticket, file", addType)
		}
		return addType, false, nil
	}

	itemType, err := DetectType(value)
	if err != nil {
		return "", false, err
	}
	return itemType, true, nil
}

func resolveContextID(s *store.Store, cfg *config.Config, addGlobal bool) (*int64, error) {
	if addGlobal {
		return nil, nil
	}

	currentCtx := cfg.CurrentContext
	if currentCtx == "" {
		currentCtx = config.GlobalContextName
	}
	c, err := s.GetContext(context.Background(), currentCtx)
	if err != nil {
		return nil, fmt.Errorf("current context %q not found", currentCtx)
	}
	return &c.ID, nil
}

func insertItem(s *store.Store, value, itemType string, contextID *int64) (*models.Item, error) {
	todoState, err := s.GetState(context.Background(), "todo")
	if err != nil {
		return nil, fmt.Errorf("default 'todo' state not found: %w", err)
	}

	itemSlug, err := nextSlug(s, slug.FromValue(value))
	if err != nil {
		return nil, err
	}

	item := &models.Item{
		Slug:      itemSlug,
		ContextID: contextID,
		Type:      itemType,
		Value:     value,
		StateID:   &todoState.ID,
	}

	if err := s.CreateItem(context.Background(), item); err != nil {
		return nil, err
	}
	return item, nil
}

func printAdded(stdout io.Writer, item *models.Item, value, itemType string, autoDetected bool, contextID *int64, cfg *config.Config) {
	typeLabel := itemType
	if autoDetected {
		typeLabel = "auto-detected as " + itemType
	}
	scope := "globally"
	if contextID != nil {
		scope = fmt.Sprintf("in context %q", cfg.CurrentContext)
	}
	fmt.Fprintf(stdout, "Added [%s] %s (%s) %s\n", item.Slug, value, typeLabel, scope)
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
