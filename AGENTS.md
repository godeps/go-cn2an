# Repository Guidelines

This guide keeps contributions aligned with the Go implementation of the `cn2an` number conversion toolkit. Review it before proposing changes or new features.

## Project Structure & Module Organization
- Core converters live in `cn2an.go`, `an2cn.go`, and `transform.go`; each file exposes a focused API for number or sentence conversion.
- Shared mappings and normalization logic are centralized in `config.go` and `normalize.go`; update these when adding new numerals or token rules.
- Unit tests reside alongside their targets (`*_test.go`), mirroring package boundaries.
- The `example/` directory contains `main.go`, a runnable showcase that should stay sync'd with the latest public API.

## Build, Test, and Development Commands
- `go build ./...` compiles every package and ensures new code respects Go 1.21 module boundaries.
- `go test ./...` runs the entire suite; prefer `go test -run TestCase ./...` when iterating on a failing scenario.
- `go run example/main.go` exercises the public API end to end and is helpful for smoke-testing transformations after refactors.
- `gofmt -w *.go` (or `go fmt ./...`) must be run prior to committing; CI assumes canonical formatting.

## Coding Style & Naming Conventions
- Follow idiomatic Go: tabs for indentation, PascalCase for exported identifiers, and camelCase for internals.
- Keep functions short and single-purpose; split helpers when a block exceeds ~30 lines or mixes responsibilities.
- Maintain consistent error messages (`fmt.Errorf("context: %w", err)`) and prefer early returns over deeply nested conditionals.

## Testing Guidelines
- Add table-driven tests in the existing `*_test.go` files; new behaviors should extend the relevant suite to avoid duplication.
- Aim to cover edge cases (negative numbers, mixed scripts, sentence punctuation) alongside the happy path.
- Run `go test -cover ./...` before opening a PR and ensure new paths keep or improve coverage.

## Commit & Pull Request Guidelines
- Commit messages follow an imperative, present-tense summary (`Add smart mode fallback`, `Fix rmb spacing bug`) with optional details in the body.
- Reference issues with `Refs #123` when applicable and keep commits logically scoped for easier review.
- PRs must include: a brief motivation, validation notes (`go test ./...`, example output), and screenshots or snippets if behavior visibly changes.
