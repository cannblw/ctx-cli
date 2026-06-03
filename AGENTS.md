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

## Architecture

### No premature interfaces

Define interfaces at the call site, not at the producer. A concrete struct (`type Store struct`) is the default. Extract an interface only when a second implementation or test mock creates the need.

### No stubs for future code

Don't add methods returning `fmt.Errorf("not implemented")`. Add them when the task that needs them lands.

## Migrations

### Embed via go:embed

Migrations are embedded with `//go:embed *.sql`, never discovered via CWD-relative path guessing. Goose reads from the embedded FS via `SetBaseFS`.

## Constants

Extract strings magic strings and numbers as constants.

## CLI output

Suppress library logs: this is a CLI tool, not a backend service. Example: use `goose.SetLogger(goose.NopLogger())` in both production and test setup.

## Examples in help text

Avoid atemporal references (version numbers, year-specific software names).

## Tests

### Naming

Every test must include `_Success` or `_Error` in its name. Prefer longer descriptive names: `TestListContexts_SuccessOneItem` not `TestListContexts_Single`.

### Ordering

Success tests first, then error tests. Simple before complex.

### Granularity

Distinct behaviors get their own test (e.g. with/without `--desc`). Fold trivial assertions (output format, ID presence) into the main Success test. Don't test framework behavior, test your own domain decisions.

### Location

Tests live next to what they test: `cmd/new_test.go`, `migrations/migrations_test.go`. Don't put migration tests in `pkg/store/` just because they share a test helper.
