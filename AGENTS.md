# AGENTS.md for ctx-cli

## Cobra CLI

For scaffolding new commands, subcommands, or flags, use `cobra-cli add` -- never write command structs or init functions from scratch.
For editing existing business logic in a command, modify its `cmd/*.go` file directly.
If cobra-cli is not installed, run `go install github.com/spf13/cobra-cli@latest` first.
