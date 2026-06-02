// Copyright © 2026 Edgar Chirivella
package main

import "github.com/cannblw/ctx-cli/cmd"

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
	}
}
