# Usage

Every command except `init` and `install` looks for `spec/` in the current directory and then in each parent directory. `xspec --version` prints the version.

## init

```sh
xspec init
```

Creates `spec/` at the git top-level (or in the current directory outside git), with the folders, a `SPEC.md` stub, `notes/glossary.md`, `ROADMAP.md`, and the kickoff `plans/001-kickoff.md` when there are no items yet. It adds `/spec/` to `.gitignore` unless git already ignores it. Existing files are kept, so running it again is safe.

## new

```sh
xspec new <type> <title> [--depth light|full] [--with 7,9] [--bug "text"]... [--dir]
```

Creates an item from its template and prints its path.

| Type | Creates |
| --- | --- |
| `feature` | `features/NNN-slug.md` |
| `issue` | `issues/NNN-slug.md` |
| `update` | `updates/NNN-slug/update.md` + `design.md` |
| `fix` | `fixes/NNN-slug.md`, or a folder with `fix.md` + `design.md` with `--dir` |
| `plan` | `plans/NNN-slug.md` |
| `note` | `notes/slug.md` (no number, no status) |
| `analysis` | `analysis/slug.md` (no number, no status) |

- `--depth light`: feature, issue and plan only. A short template: Goal / Tasks for a feature, Current / Expected / Tasks for an issue, Context / Decision / Steps for a plan. `full` is the default.
- `--with 7,9`: items to include. An update takes features and issues; a fix takes issues. Run in a terminal without `--with`, `new update` and `new fix` show a numbered list and ask which to include.
- `--bug "text"`: fix only, repeatable. Adds a line under Small bugs.
- `--dir`: fix only. Makes a folder with `design.md`.

## list

```sh
xspec list [--all]
```

Prints active numbered items: number, type, status, the count of open `[?: …]` markers, checked/total task boxes, and title. `--all` adds archived items.

## next

```sh
xspec next
```

Prints one line: the next unchecked task of the lowest `doing` item; else an item to close; else an approved item to start; else an item waiting for answers or approval; else `nothing to do`.

## questions

```sh
xspec questions
```

Prints every open `[?: …]` marker in active items as `file:line  question`, grouped by item. Markers inside backticks or code blocks are examples and do not count.

## approve

```sh
xspec approve <num|slug>
```

The owner's approval. Sets `status: approved <hash>` on the item and on every open item it includes, and refuses while any of them has a marker. The hash covers the content without the status line, the Change line and task-box state: ticking tasks keeps the approval, any other edit voids it until the owner approves again. Agents must not run it.

## start

```sh
xspec start <num|slug>
```

Sets `status: doing` on the item and on every active item it includes. Each must be approved and unchanged since, or already doing.

## done

```sh
xspec done <num|slug> [--change <commit|PR|path>] [--no-verify]
```

For the item and every active item it includes, it:

1. refuses while a `[?: …]` marker or an unchecked task box remains, or an item was never approved;
2. runs every `Verify: <command>` line once, in the project root, and refuses if one fails (`--no-verify` skips them);
3. ticks the Includes boxes;
4. sets `status: done` and writes `Change: <commit> <YYYY-MM-DD>` on line 2 (the commit from `--change`, else the short git HEAD, else only the date);
5. moves the file or folder to `spec/archive/<type folder>/`.

## roadmap

```sh
xspec roadmap
```

Rewrites and prints `spec/ROADMAP.md`. Every `new`, `approve`, `start` and `done` rewrites it too. Sections: Now (doing), Next (approved, then open, with question counts), History (archived, newest first, with date and commit). Items included in an update or fix are listed under it.

## context

```sh
xspec context <num|slug>
```

Prints everything an agent needs for one item: the item (and `design.md`), the items it includes, the Principles section of `SPEC.md`, the kickoff's Stack table and `notes/glossary.md`, each under a `=== path ===` header. Paste it into an agent that cannot read `spec/`.

## check

```sh
xspec check
```

Prints every problem and exits 1 when there is one: `spec/` not ignored by git, a bad status line, an approved item without a hash or edited after approval, a duplicate number, `doing` with open markers, `done` outside the archive or an archived item that is not done, an Includes line pointing at a missing number, and a requirement (`- R1: …`) that no task box references. It fits a pre-commit hook.

## install

```sh
xspec install project|global [--force]
```

Writes the bundled execspec skill into `.claude/skills/execspec/` and `.agents/skills/execspec/`, under the project's git top-level (`project`) or your home directory (`global`). It refuses to touch an existing skill directory unless `--force` is set; `--force` replaces it whole.
