<div align="center">

# execspec

**Spec-driven development for coding with AI agents, in a private `spec/` folder.**

[![Go 1.27+](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
![Local first](https://img.shields.io/badge/local--first-no%20LLM%20calls%2C%20no%20network-informational)

English · [Русский](README_RU.md)

</div>

> [!WARNING]
> The project is in early development and may change significantly: commands, flags, the `spec/` file format and the agent skill can change between versions without backward compatibility. Pin a version and check the release notes before upgrading.

`xspec` is a small CLI that keeps your intent next to the code: numbered changes and updates as plain Markdown in a git-ignored `spec/` folder, next to a short project brief (`SPEC.md`). You and your AI agent read and write the same files. The CLI creates them from templates, closes them with `open → done` and refuses to close while questions or tasks are still open. It ships an agent skill that teaches Claude Code, Codex and other agents the same workflow.

## Features

- **Private by default.** `spec/` holds intent: plans, opinions, open questions. `xspec init` keeps it out of git.
- **Numbered items.** One counter for both types: `xspec done 7`, `--with 7,9`.
- **Testable requirements.** Add EARS requirements (`WHEN … THE SYSTEM SHALL …`) when behavior needs them, with tasks that name them (`(R1)`). A bug records current, expected *and unchanged* behavior.
- **The agent works, you answer.** The agent shows a short plan in chat and starts. It stops only at an open `[?: question]` marker or an *Ask first* area of `SPEC.md`. Unchecked tasks and failing `Verify:` commands block `done`.
- **Right-sized specs.** A one-sentence task needs no file. A *change* is one piece of work, an *update* is a batch; both start at Goal / Tasks / Verify and grow.
- **A brief the agent drafts.** `SPEC.md` holds Context, Audience, Principles, Design, Stack and Ask first. The agent writes the draft from your repository and a few questions.
- **Always know what is next.** `xspec next`, `xspec questions`, and a generated `ROADMAP.md` with history.
- **Clean history.** `done` writes the commit into the file and moves it into `spec/archive/`; `--decision` keeps one line of reasoning in `notes/decisions.md`.
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
xspec init                                  # spec/, the SPEC.md brief, a .gitignore line
xspec new change "Login form"               # spec/changes/001-login-form.md
xspec new change "Crash on save" --bug      # spec/changes/002-crash-on-save.md
xspec new update v0.2 --with 1,2            # spec/updates/003-v0-2.md (--dir adds design.md)
xspec questions                             # what waits for your answer
xspec next                                  # the next task to do
xspec done 3                                # refused while tasks are unchecked or Verify fails
xspec check                                 # exit 1 on any problem
```

### What a session looks like

```text
$ xspec list
NUM  TYPE    STATUS  [?]  TASKS  TITLE
001  change  open    0    1/2    Login form
002  change  open    0    0/2    Crash on save
003  update  open    0    0/1    v0.2

$ xspec done 3 --decision "sessions over JWT: simpler to revoke"
verify: go test ./...
003 done   v0.2 -> spec/archive/updates/003-v0-2.md
001 done   Login form -> spec/archive/changes/001-login-form.md
002 done   Crash on save -> spec/archive/changes/002-crash-on-save.md
```

A new change file looks like this. The agent fills it in, shows a short plan, and writes the code:

```markdown
status: open
# 001 Login form

## Goal
[?: what changes, and for whom]

## Tasks
- [ ] T1 [?: task]

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
```

## Documentation

| Page | What is in it |
| --- | --- |
| [Usage](docs/user/usage.md) | Every command and flag |
| [Workflow](docs/user/workflow.md) | The spec layout, item types, the brief and the work cycle |
| [Architecture](docs/dev/architecture.md) | Internals for contributors and coding agents; start with [AGENTS.md](AGENTS.md) |

## Development

```sh
make          # fmt, vet, test, build → build/bin/xspec
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
