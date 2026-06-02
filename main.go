// Copyright © 2026 Edgar Chirivella
package main

import "github.com/cannblw/ctx-cli/cmd"

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		// cobra already prints the error, just exit non-zero
		// (in future tasks we'll wire real error handling)
	}
}
