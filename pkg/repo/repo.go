package repo

import (
	"context"

	"github.com/cannblw/ctx-cli/pkg/models"
)

type Repository interface {
	CreateContext(ctx context.Context, name, description string) (*models.Context, error)
	GetContext(ctx context.Context, name string) (*models.Context, error)
	ListContexts(ctx context.Context) ([]models.Context, error)
	Close() error
}
