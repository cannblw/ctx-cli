package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/models"
	"github.com/cannblw/ctx-cli/pkg/slug"
	"github.com/cannblw/ctx-cli/pkg/store"
)

var validTypeNames = []string{"link", "pr", "ticket", "file"}

var validTypes = buildTypeMap()

var validTypesStr = strings.Join(validTypeNames, ", ")

func buildTypeMap() map[string]struct{} {
	m := make(map[string]struct{}, len(validTypeNames))
	for _, t := range validTypeNames {
		m[t] = struct{}{}
	}
	return m
}

// NewAddCmd creates the `ctx add` command.
func NewAddCmd(s *store.Store, cfg *config.Config, stdout io.Writer) *cobra.Command {
	var (
		itemType string
		isGlobal bool
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

			itemType, autoDetected, err := resolveItemType(value, itemType)
			if err != nil {
				return err
			}

			contextID, err := ResolveContextID(s, cfg, isGlobal)
			if err != nil {
				return err
			}

			if itemType == "file" {
				value, err = copyFileToCtxDir(value, resolveContextName(cfg, isGlobal))
				if err != nil {
					return err
				}
			}

			item, err := insertItem(s, value, itemType, contextID, stateName(cfg))
			if err != nil {
				return err
			}

			printAdded(stdout, item, value, itemType, autoDetected, contextID, cfg)
			return nil
		},
	}

	cmd.Flags().StringVarP(&itemType, "type", "t", "", "item type: "+validTypesStr)
	cmd.Flags().BoolVarP(&isGlobal, "global", "g", false, "add item globally (no context)")
	return cmd
}

// DetectType auto-detects the item type from the value.
func DetectType(value string) (string, error) {
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return "link", nil
	}
	path := value
	if strings.HasPrefix(value, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = home + value[1:]
		}
	}
	if _, err := os.Stat(path); err == nil {
		return "file", nil
	}
	return "", fmt.Errorf("could not detect type for %q, use --type to specify one of "+validTypesStr, value)
}

func resolveItemType(value, itemType string) (string, bool, error) {
	if itemType != "" {
		if _, ok := validTypes[itemType]; !ok {
			return "", false, fmt.Errorf("invalid type %q: must be one of "+validTypesStr, itemType)
		}
		return itemType, false, nil
	}

	itemType, err := DetectType(value)
	if err != nil {
		return "", false, err
	}
	return itemType, true, nil
}

func copyFileToCtxDir(src, ctxName string) (string, error) {
	destDir := config.ContextDir(ctxName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("could not create context dir: %w", err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer srcFile.Close()

	dest := filepath.Join(destDir, filepath.Base(src))
	dstFile, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("could not create file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", fmt.Errorf("could not copy file: %w", err)
	}
	return dest, nil
}

func stateName(cfg *config.Config) string {
	if cfg.DefaultItemState != "" {
		return cfg.DefaultItemState
	}
	return config.DefaultItemState
}

func insertItem(s *store.Store, value, itemType string, contextID *int64, stateName string) (*models.Item, error) {
	state, err := s.GetState(context.Background(), stateName)
	if err != nil {
		return nil, fmt.Errorf("default state %q not found: %w", stateName, err)
	}

	itemSlug, err := resolveSlug(s, value)
	if err != nil {
		return nil, err
	}

	item := &models.Item{
		Slug:      itemSlug,
		ContextID: contextID,
		Type:      itemType,
		Value:     value,
		StateID:   &state.ID,
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

func resolveSlug(s *store.Store, value string) (string, error) {
	base := slug.FromValue(value)
	count, err := s.CountSlugsByPrefix(context.Background(), base)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return base, nil
	}
	return fmt.Sprintf("%s-%d", base, count+1), nil
}

func resolveContextName(cfg *config.Config, isGlobal bool) string {
	if isGlobal {
		return ""
	}
	return cfg.CurrentContext
}
