package store

import (
	"context"
	"fmt"

	"github.com/cannblw/ctx-cli/pkg/models"
)

func (s *Store) GetState(ctx context.Context, name string) (*models.State, error) {
	state := &models.State{}
	err := s.db.NewSelect().Model(state).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get state %q: %w", name, err)
	}
	return state, nil
}
