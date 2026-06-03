package store

import (
	"context"
	"fmt"

	"github.com/cannblw/ctx-cli/pkg/models"
)

func (s *Store) CreateContext(ctx context.Context, name, description string) (*models.Context, error) {
	c := &models.Context{
		Name:        name,
		Description: description,
	}
	_, err := s.db.NewInsert().Model(c).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("create context: %w", err)
	}
	return c, nil
}

func (s *Store) GetContext(ctx context.Context, name string) (*models.Context, error) {
	c := new(models.Context)
	err := s.db.NewSelect().Model(c).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get context %q: %w", name, err)
	}
	return c, nil
}

func (s *Store) ListContexts(ctx context.Context) ([]models.Context, error) {
	var contexts []models.Context
	err := s.db.NewSelect().Model(&contexts).OrderExpr("created_at ASC").Scan(ctx)
	return contexts, err
}
