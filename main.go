package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cannblw/ctx-cli/cmd"
	"github.com/cannblw/ctx-cli/pkg/store"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot find home directory:", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(home, ".ctx", "ctx.db")

	r, err := store.NewBunStore(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot open database:", err)
		os.Exit(1)
	}
	defer r.Close()

	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewNewCmd(r, os.Stdout))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
