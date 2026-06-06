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
		return nil, fmt.Errorf("could not insert context: %w", err)
	}
	return c, nil
}

func (s *Store) GetContext(ctx context.Context, name string) (*models.Context, error) {
	c := &models.Context{}
	err := s.db.NewSelect().Model(c).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get context %q: %w", name, err)
	}
	return c, nil
}

func (s *Store) ListContexts(ctx context.Context) ([]models.Context, error) {
	var contexts []models.Context
	err := s.db.NewSelect().Model(&contexts).OrderExpr("created_at ASC").Scan(ctx)
	return contexts, err
}

func (s *Store) RenameContext(ctx context.Context, oldName, newName string) (*models.Context, error) {
	_, err := s.GetContext(ctx, newName)
	if err == nil {
		return nil, fmt.Errorf("could not rename context: %q already exists", newName)
	}

	res, err := s.db.NewUpdate().
		Model((*models.Context)(nil)).
		Set("name = ?", newName).
		Set("updated_at = datetime('now')").
		Where("name = ?", oldName).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not rename context: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("could not rename context %q: database update matched no rows", oldName)
	}

	return s.GetContext(ctx, newName)
}

func (s *Store) DeleteContext(ctx context.Context, name string) error {
	_, err := s.GetContext(ctx, name)
	if err != nil {
		return fmt.Errorf("context %q not found", name)
	}

	_, err = s.db.NewDelete().
		Model((*models.Context)(nil)).
		Where("name = ?", name).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("could not delete context %q: %w", name, err)
	}
	return nil
}
