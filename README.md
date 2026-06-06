# ctx

<p align="center">
  <em>Harness context switching. Keep your flow.</em>
</p>

`ctx` helps you track and switch between multiple workstreams so you can stay focused when juggling agentic flows, side projects, or anything that pulls your attention in different directions.

## Why ctx?

When you're deep in one problem but need to jump to another — whether it's a different codebase, a code review, or a conversation with an AI agent — it's easy to lose your mental state. `ctx` gives you a lightweight way to save where you are, switch contexts, and pick up right where you left off.

## Install

```sh
go install github.com/cannblw/ctx-cli@latest
```

Or build from source:

```sh
git clone https://github.com/cannblw/ctx-cli.git
cd ctx-cli
make build
```

The binary lands at `bin/ctx`.

## Quick start

```sh
ctx new my-context --desc "Refactoring the auth module"
ctx contexts
```

## Usage

```sh
# Create contexts
ctx new fix-auth --description "Fix the login redirect bug"
ctx new migrate-db --description "Database migration"

# List contexts
ctx contexts          # or: ctx c, ctx ctx

# Switch to a context
ctx switch fix-auth   # or: ctx s fix-auth

# `ctx contexts <name>` also switches
ctx contexts fix-auth

# Rename a context
ctx rename fix-auth fix-auth-bug
# or: ctx rn fix-auth fix-auth-bug

# Delete a context
ctx rm --context old-project
# or: ctx remove --context old-project
# Add --force / -f to skip confirmation

# Add items to a context

# Switch to a context first
ctx switch my-project

# Add a link
ctx add https://github.com/org/repo/pull/42

# Add a link with explicit type
ctx add https://jira.company.com/PROJ-123 --type ticket

# Add a file
ctx add ~/notes/todo.md

# Add globally (no context needed)
ctx add https://example.com --global

# Aliases
ctx a https://example.com
```

## License

MIT © 2026 Edgar Chirivella
