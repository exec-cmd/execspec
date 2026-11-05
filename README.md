<div align="center">

# execspec

**Spec-driven development for coding with AI agents, in a private `spec/` folder.**

[![Go 1.27+](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
![Local first](https://img.shields.io/badge/local--first-no%20LLM%20calls%2C%20no%20network-informational)

English · [Русский](README_RU.md)

</div>

`xspec` is a small CLI that keeps your intent next to the code: numbered features, bugs, updates, fixes and plans as plain Markdown in a git-ignored `spec/` folder. You and your AI agent read and write the same files. The CLI creates them from templates, moves them through `open → doing → done` and refuses to move on while questions are still open. It ships an agent skill that teaches Claude Code, Codex and other agents the same workflow.

## Features

- **Private by default.** `spec/` holds intent: plans, opinions, open questions. `xspec init` keeps it out of git.
- **Numbered items.** One counter for all types: `xspec start 7`, `xspec done 7`, `--with 7,9`.
- **Testable requirements.** Features use EARS requirements (`WHEN … THE SYSTEM SHALL …`) and tasks that name them (`(R1)`). Bugs record current, expected *and unchanged* behavior.
- **You approve, the agent builds.** `[?: question]` markers block approval; `xspec approve` is yours alone, and any edit afterwards voids it. Unchecked tasks and failing `Verify:` commands block `done`.
- **Right-sized specs.** A kickoff for the project, `--depth light` for small work, one update with parts when you hand over a batch.
- **Always know what is next.** `xspec next`, `xspec questions`, and a generated `ROADMAP.md` with history.
- **Batches.** An *update* bundles features and issues, a *fix* bundles issues and small bugs. Closing the batch closes everything in it.
- **Clean history.** `done` writes the commit into the file and moves it into `spec/archive/`.
- **Agent skill included.** `xspec install global` installs it for Claude Code (`~/.claude/skills`) and other agents (`~/.agents/skills`).
- **No LLM inside.** No keys, no network, one static binary.

## Install

Download the binary for your system from [Releases](https://github.com/exec-cmd/execspec/releases), rename it to `xspec`, and put it on your `PATH`:

```sh
chmod +x xspec-linux-amd64 && mv xspec-linux-amd64 ~/.local/bin/xspec
```

Or with Go 1.27+:

```sh
go install github.com/exec-cmd/execspec/src/xspec@latest
# or from a checkout:
make build    # → build/bin/xspec
```

Then install the agent skill:

```sh
xspec install global     # ~/.claude/skills/execspec and ~/.agents/skills/execspec
```

## Quick start

```sh
xspec init                                  # spec/, the kickoff, a .gitignore line
xspec new feature "Login form"              # spec/features/002-login-form.md
xspec new issue "Crash on save"             # spec/issues/003-crash-on-save.md
xspec new update v0.2 --with 2,3            # spec/updates/004-v0-2/{update,design}.md
xspec questions                             # what waits for your answer
xspec approve 4                             # yours: refused while [?: …] markers remain
xspec start 4
xspec next                                  # the next task to do
xspec done 4                                # refused while tasks are unchecked or Verify fails
xspec check                                 # exit 1 on any problem
```

### What a session looks like

```text
$ xspec list
NUM  TYPE     STATUS  [?]  TASKS  TITLE
002  feature  open    0    0/2    Login form
003  issue    open    0    0/2    Crash on save
004  update   open    0    0/1    v0.2

$ xspec done 4
verify: go test ./...
004 done   v0.2 -> spec/archive/updates/004-v0-2
002 done   Login form -> spec/archive/features/002-login-form.md
003 done   Crash on save -> spec/archive/issues/003-crash-on-save.md
```

A new feature file looks like this. The agent fills it in, you approve it, and then the code gets written:

```markdown
status: open
# 002 Login form

## Why
[?: the problem, and who gets what when it is solved]

## Requirements
- R1: WHEN [?: condition] THE SYSTEM SHALL [?: behavior]

## Design
[?: the approach and the code it touches]

## Tasks
- [ ] T1 [?: task] (R1)

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
```

## Documentation

| Page | What is in it |
| --- | --- |
| [Usage](docs/user/usage.md) | Every command and flag |
| [Workflow](docs/user/workflow.md) | The spec layout, item types, templates and the work cycle |
| [Architecture](docs/dev/architecture.md) | Internals for contributors and coding agents; start with [AGENTS.md](AGENTS.md) |

## Development

```sh
make          # fmt, vet, test, build → build/bin/xspec
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
