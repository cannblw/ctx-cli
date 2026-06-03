package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/store"
)

const (
	DBDir  = ".ctx"
	DBFile = "ctx.db"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not find home directory:", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(home, DBDir, DBFile)

	s, err := store.NewStore(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open database:", err)
		os.Exit(1)
	}
	defer s.Close()

	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewNewCmd(s, os.Stdout))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
