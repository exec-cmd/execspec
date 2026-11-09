# Architecture

## Packages

| Package | Role |
| --- | --- |
| `src/xspec` | `main`: calls `cli.Execute`. |
| `src/cli` | Cobra commands, flags, the terminal pick list for `new update`, output. No file logic. |
| `src/spec` | Everything about `spec/`: layout and `Init` (spec.go), item scan and parse, numbering, `New`, `Done` (item.go), templates (templates.go), `Check` (check.go), `Roadmap`, `Next`, `Questions`, `Context` (view.go). No Cobra, no printing. |
| `src/skill` | The agent skill (source of truth) embedded with `go:embed`, and `Install`. |

## Items

`spec.Scan` reads the status-type folders (`changes`, `updates`) under `spec/` and `spec/archive/`. An entry counts when its name matches `NNN-slug.md` (a file) or `NNN-slug` (a folder whose main file is `<type>.md`). Everything else is ignored.

Parsing is line-based (`Item.parse`):

| What | Rule |
| --- | --- |
| status | line 1 `status: open` or `status: done`; `Item.Started` is an open item with a checked task |
| change | line 2 `Change: <what> <YYYY-MM-DD>`; the date is optional (v0.2 wrote none) |
| title | first `# ` heading, number prefix removed |
| markers | `[?:…]` with line numbers, outside code spans and fenced blocks |
| verify | `Verify: <command>` lines |
| includes | `- [ ] NNN change: …`; not counted as tasks |
| tasks | every other `- [ ]` / `- [x]` line |
| requirements | `- R<n>:` lines, checked against `R<n>` mentions in task lines |

## Numbering

The next number is the highest number in `Scan` + 1, so archived items count and numbers are never reused. Numbers are zero-padded to three digits and grow past 999 without change.

## done

`Done` resolves the ref (number or slug), then builds a *group*: the item plus every active item its Includes lines name. Validation runs over the whole group before any file is written. There is no approval step and no content hash: the stops are open `[?: …]` markers (answered by the owner in chat) and the Ask first section of `SPEC.md` (read by the skill).

`Done` needs all boxes checked and no markers; then runs each distinct Verify command through `DoneOptions.Run` (the CLI passes `Shell`, tests pass a stub), appends the `--decision` line to `notes/decisions.md` (`appendDecision`), ticks the Includes boxes, rewrites status and Change, and renames each member into `spec/archive/<type folder>/`.

Progress is derived, not stored: `Roadmap` puts an open item with a checked task under Now, `Next` names its first unchecked task, and `Check` reports a started item that still has markers.

Every mutating operation ends with `WriteRoadmap`. `Roadmap` lists an item under its parent when an item on the same side of the archive includes it.

## Git

The CLI shells out to `git` only for three reads: `rev-parse --show-toplevel` (project root), `check-ignore -q spec/` (ignored or not), and `rev-parse --short HEAD` (Change line). Outside git, init checks `.gitignore` lines instead, and `check` skips the ignore rule.

## Skill

`src/skill/SKILL.md` and `src/skill/references/` are the skill as shipped. Edit them there; `xspec install` copies the embedded version. Keep the skill and `src/spec/templates.go` in step: the skill describes the same format the CLI writes. Maintainer-only files stay out of `src/skill/`, because `go:embed` ships the whole `references/` folder; the acceptance scenarios live in `docs/dev/skill-acceptance.md`.

`references/polish.md` holds the improve and digest rules for every spec file. The digest is never stored: a stored summary goes stale beside the item, and a file named like an item would be scanned as one (`nameRe`).

## Version

`cli.version` is set with `-ldflags -X` by the Makefile (`git describe`) and the release workflow (the tag); a plain `go build` reports `dev`.

## Tests

`src/spec/spec_test.go` runs every operation in `t.TempDir()`, outside git.
