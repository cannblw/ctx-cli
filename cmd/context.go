package cmd

import (
	"context"
	"fmt"

	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/store"
)

// ResolveContextName returns "" for global items, otherwise the current context name.
func ResolveContextName(cfg *config.Config, isGlobal bool) string {
	if isGlobal {
		return ""
	}
	return cfg.CurrentContext
}

func ResolveContextID(s *store.Store, cfg *config.Config, isGlobal bool) (*int64, error) {
	if isGlobal || cfg.CurrentContext == "" {
		return nil, nil
	}

	c, err := s.GetContext(context.Background(), cfg.CurrentContext)
	if err != nil {
		return nil, fmt.Errorf("current context %q not found", cfg.CurrentContext)
	}
	return &c.ID, nil
}
