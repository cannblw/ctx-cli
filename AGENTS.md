# AGENTS.md for ctx-cli

## Cobra CLI

Scaffold new commands with `cobra-cli add` — don't write command structs or init functions from scratch.
Edit business logic directly in `cmd/*.go`.
If cobra-cli is missing: `go install github.com/spf13/cobra-cli@latest`.

## Dependency injection by constructor — no setters, no globals

**Never** use package-level globals or setters (`SetRepo()`, `SetConfig()`, `getRepo()`) to pass dependencies.

**Always** use constructor functions.

## Comments

**Minimal** comments. Never explain *what* the code does — the code itself should make that obvious. Only explain *why* when the rationale is genuinely non-obvious. If the intent is clear from reading the code, skip the comment.

**Doc comments for all public things.** Every exported function, method, type, constant, and variable must have a doc comment starting with the name of the thing it describes (Go convention).
