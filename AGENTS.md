# AGENTS.md for ctx-cli

## Cobra CLI

Scaffold new commands with `cobra-cli add` — don't write command structs or init functions from scratch.
Edit business logic directly in `cmd/*.go`.
If cobra-cli is missing: `go install github.com/spf13/cobra-cli@latest`.

## Dependency injection by constructor — no setters, no globals

**Never** use package-level globals or setters (`SetRepo()`, `SetConfig()`, `getRepo()`) to pass dependencies.

**Always** use constructor functions.
