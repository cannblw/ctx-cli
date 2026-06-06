package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/config"
	"github.com/cannblw/ctx-cli/pkg/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not load config:", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(config.Dir(), "ctx.db")

	s, err := store.NewStore(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open database:", err)
		os.Exit(1)
	}
	defer s.Close()

	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(
		cmd.NewNewCmd(s, os.Stdout),
		cmd.NewContextsCmd(s, cfg, os.Stdout),
		cmd.NewSwitchCmd(s, cfg, os.Stdout),
		cmd.NewRenameCmd(s, cfg, os.Stdout),
		cmd.NewRmCmd(s, cfg, os.Stdin, os.Stdout),
		cmd.NewAddCmd(s, cfg, os.Stdout),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
