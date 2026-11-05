## Project overview

`execspec` ships `xspec`, a Go 1.27 single-binary CLI (`github.com/exec-cmd/execspec`) using Cobra. It manages a private, git-ignored `spec/` directory for spec-driven development: numbered features, issues, updates, fixes and plans in Markdown, with statuses, gates and an archive. It also embeds and installs the `execspec` agent skill. It makes no LLM or network calls.

## Repository layout

- `src/xspec/main.go`: entry point.
- `src/cli/`: Cobra commands, flags, prompts and output only.
- `src/spec/`: layout, templates, item parsing, numbering, start/done/archive, check rules; independent of Cobra.
- `src/spec/spec_test.go`: tests in temp directories.
- `src/skill/`: the agent skill (source of truth) and its installer.
- `docs/user/`: usage and workflow for users. `docs/dev/`: internals; read `docs/dev/architecture.md` before changing `src/spec`.
- `.github/workflows/`: CI verification and tagged cross-platform releases.
- `Makefile`: fmt, vet, test, build.

## Usage and verification

```sh
go run ./src/xspec --help
go test ./...
go vet ./...
make            # fmt, vet, test, then build/bin/xspec
```

## Change rules

- Keep CLI code in `src/cli` and spec logic in `src/spec`.
- A change to the file format touches three places together: `src/spec/templates.go` (and the parsers), the skill in `src/skill/`, and `docs/user/workflow.md`.
- Run `gofmt` (or `make fmt`) after Go changes and add focused tests for new behavior.
- Report only checks that were actually run.
