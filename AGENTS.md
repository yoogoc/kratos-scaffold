# Repository Guidelines

## Project Structure & Module Organization
- cmd/: Cobra CLI commands (e.g., `root.go`, `proto.go`, `biz.go`, `data.go`, `service.go`, `generate.go`), plus tests and `testdata/`.
- generator/: Code generation logic and templates (`tmpl/`, `wire/`) for proto, biz, data, service.
- project_generator/: Project scaffolding logic and embedded resources (`resources/`, `project.go`, `normal.go`, `bff.go`).
- pkg/: Supporting libraries (`field/`, `util/`, `cli/`).
- internal/: Internal utilities (e.g., `merge/`).
- main.go: CLI entrypoint. Module: `github.com/yoogoc/kratos-scaffold`.

## Build, Test, and Development Commands
- Build: `go build -o bin/kratos-scaffold .` — compile the CLI locally.
- Install: `go install github.com/yoogoc/kratos-scaffold@latest` — install into `$GOBIN`.
- Run locally: `go run . --help` or `go run . proto ...` — execute from source.
- Test: `go test ./...` (with coverage: `go test ./... -cover`).
- Lint/Format: `go vet ./...` and `go fmt ./...` before pushing.

## Coding Style & Naming Conventions
- Formatting: Use standard Go formatting (`go fmt`); keep imports tidy; run `go vet`.
- Naming: packages lower-case; exported identifiers in CamelCase; files lower_snake_case where helpful.
- CLI design: keep subcommands concise (e.g., `proto`, `biz`, `data`, `service`, `new`, `g`); prefer clear, short flags.

## Testing Guidelines
- Framework: standard `testing` package; tests in `*_test.go`, functions `TestXxx`.
- Patterns: prefer table-driven tests; place fixtures under `testdata/` (see `cmd/testdata/`).
- Scope: add/extend tests when changing generation logic in `generator/` or command behavior in `cmd/`.

## Commit & Pull Request Guidelines
- Commits: small, focused, imperative mood (e.g., "generator: fix time type", "cmd: add json,array type"). Use a short subject and optional body for context.
- PRs: include a clear description, linked issues, and examples of expected CLI usage (input command and generated paths). Paste `go test ./...` output when relevant.

## Security & Configuration Tips
- Toolchain: use a recent Go matching `go.mod`. Set `GOBIN`/`GOPATH` appropriately.
- Safety: generation can overwrite files; run on a clean branch and review diffs. For proto and service paths, pass explicit `-o` or namespace flags to avoid surprises.
