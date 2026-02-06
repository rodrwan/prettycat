# Repository Guidelines

## Project Structure & Module Organization
This repository contains a Go CLI app, `prettycat`.

- `cmd/prettycat/main.go`: CLI entrypoint and flag parsing.
- `internal/app`: top-level orchestration (input, render, output, exit codes).
- `internal/input`: file/stdin loading and error handling.
- `internal/render`: Markdown/code/plain rendering logic and detection.
- `internal/pager`: interactive terminal pager.
- `internal/style`: shared ANSI/header styling helpers.
- `testdata/`: sample files used for manual checks and examples.
- `README.md`: user-facing usage docs.
- `Makefile`: common developer tasks.

Prefer adding new implementation code under `internal/` and keep `cmd/` thin.

## Build, Test, and Development Commands
Use the `Makefile` targets:

- `make build`: build binary at `bin/prettycat`.
- `make run ARGS="testdata/sample.md"`: run locally with arguments.
- `make test`: run all tests (`go test ./...`).
- `make testv`: verbose test output.
- `make fmt`: format all Go files with `gofmt`.
- `make fmt-check`: fail if formatting is needed.
- `make tidy`: sync module deps (`go mod tidy`).
- `make check`: format check + tests + build.
- `make clean`: remove build artifacts.

## Coding Style & Naming Conventions
- Language: Go (idiomatic style).
- Formatting: always run `gofmt` (tabs, standard import grouping).
- Packages: short, lowercase names (`render`, `pager`).
- Exported identifiers: `CamelCase`; internal helpers: `camelCase`.
- Keep functions focused and side effects explicit (especially in `internal/app`).

## Testing Guidelines
- Framework: Go `testing` package.
- Test files end with `_test.go`.
- Test names: `TestXxx` with behavior-oriented cases (table tests where useful).
- Add regression tests for rendering/pager bugs (ANSI handling, no-color behavior, error codes).
- Run `make test` before opening a PR.

## Commit & Pull Request Guidelines
No strict historical convention is enforced here yet; use clear, imperative commits.

- Recommended format: `component: concise change` (example: `pager: fix raw mode line wrapping`).
- Keep commits scoped (one logical change per commit).
- PRs should include:
  - What changed and why.
  - How to validate (`make test`, manual run examples).
  - Terminal screenshots/gifs for pager or rendering UX changes.

## Security & Configuration Notes
- Do not hardcode secrets or machine-specific paths.
- Preserve default terminal-safe behavior (`--no-color`, non-TTY output handling).
