# Architecture

## Packages

| Package | Role |
| --- | --- |
| `src/xspec` | `main`: calls `cli.Execute`. |
| `src/cli` | Cobra commands, flags, the terminal pick list for `new update/fix`, output. No file logic. |
| `src/spec` | Everything about `spec/`: layout and `Init`, item scan and parse, numbering, templates, `New`, `Start`, `Done`, `Check`. No Cobra, no printing. |
| `src/skill` | The agent skill (source of truth) embedded with `go:embed`, and `Install`. |

## Items

`spec.Scan` reads the status-type folders (`features`, `issues`, `updates`, `fixes`, `plans`) under `spec/` and `spec/archive/`. An entry counts when its name matches `NNN-slug.md` (a file) or `NNN-slug` (a folder whose main file is `<type>.md`). Everything else is ignored.

Parsing is line-based (`Item.parse`):

| What | Rule |
| --- | --- |
| status | line 1 `status: X` |
| change | line 2 `Change: X` |
| title | first `# ` heading, number prefix removed |
| markers | `[?:…]` anywhere |
| includes | `- [ ] NNN feature|issue: …`; not counted as tasks |
| tasks | every other `- [ ]` / `- [x]` line |
| requirements | `- R<n>:` lines, checked against `R<n>` mentions in task lines |

## Numbering

The next number is the highest number in `Scan` + 1, so archived items count and numbers are never reused. Numbers are zero-padded to three digits and grow past 999 without change.

## start and done

Both resolve the ref (number or slug), then build a *group*: the item plus every active item its Includes lines name. Validation runs over the whole group before any file is written. `Done` then ticks the Includes boxes, rewrites the status and Change lines of each member, and renames each into `spec/archive/<type folder>/`.

## Git

The CLI shells out to `git` only for three reads: `rev-parse --show-toplevel` (project root), `check-ignore -q spec/` (ignored or not), and `rev-parse --short HEAD` (Change line). Outside git, init checks `.gitignore` lines instead, and `check` skips the ignore rule.

## Skill

`src/skill/SKILL.md` and `src/skill/references/` are the skill as shipped. Edit them there; `xspec install` copies the embedded version. Keep the skill and `src/spec/templates.go` in step: the skill describes the same format the CLI writes.

## Tests

`src/spec/spec_test.go` runs every operation in `t.TempDir()`, outside git.
