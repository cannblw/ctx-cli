package cmd

import (
	"context"
	"fmt"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/store"
)

// resolveContextID returns nil for global items, otherwise the current context's DB ID.
func resolveContextID(s *store.Store, cfg *config.Config, isGlobal bool) (*int64, error) {
	if isGlobal {
		return nil, nil
	}

	currentCtx := cfg.CurrentContextName()
	c, err := s.GetContext(context.Background(), currentCtx)
	if err != nil {
		return nil, fmt.Errorf("current context %q not found", currentCtx)
	}
	return &c.ID, nil
}
