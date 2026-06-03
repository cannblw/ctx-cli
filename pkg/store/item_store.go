package store

import (
	"context"
	"fmt"

	"github.com/cannblw/ctx-cli/pkg/models"
)

func (s *Store) CreateItem(ctx context.Context, item *models.Item) error {
	_, err := s.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		return fmt.Errorf("could not insert item: %w", err)
	}
	return nil
}

func (s *Store) GetItem(ctx context.Context, identifier string) (*models.Item, error) {
	item := new(models.Item)
	err := s.db.NewSelect().Model(item).Where("slug = ?", identifier).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get item %q: %w", identifier, err)
	}
	return item, nil
}
