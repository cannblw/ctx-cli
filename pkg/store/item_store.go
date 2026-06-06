package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cannblw/ctx-cli/pkg/models"
)

func (s *Store) CreateItem(ctx context.Context, item *models.Item) error {
	_, err := s.GetItem(ctx, item.Slug)
	if err == nil {
		return fmt.Errorf("item %q already exists", item.Slug)
	}

	_, err = s.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		return fmt.Errorf("could not insert item: %w", err)
	}
	return nil
}

func (s *Store) GetItem(ctx context.Context, identifier string) (*models.Item, error) {
	item := &models.Item{}
	err := s.db.NewSelect().Model(item).Where("slug = ?", identifier).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get item %q: %w", identifier, err)
	}
	return item, nil
}

// CountSlugsByPrefix returns the number of items whose slug matches prefix or prefix-<number>.
func (s *Store) CountSlugsByPrefix(ctx context.Context, prefix string) (int, error) {
	var slugs []string
	err := s.db.NewSelect().
		Model((*models.Item)(nil)).
		Column("slug").
		Where("slug = ? OR slug LIKE ?", prefix, prefix+"-%").
		Scan(ctx, &slugs)
	if err != nil {
		return 0, fmt.Errorf("could not count slugs by prefix %q: %w", prefix, err)
	}

	count := 0
	for _, slug := range slugs {
		if slug == prefix {
			count++
			continue
		}
		rest := strings.TrimPrefix(slug, prefix+"-")
		if _, err := strconv.Atoi(rest); err == nil {
			count++
		}
	}
	return count, nil
}
