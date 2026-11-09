# Usage

Every command except `init` and `install` looks for `spec/` in the current directory and then in each parent directory. `xspec --version` prints the version.

## init

```sh
xspec init
```

Creates `spec/` at the git top-level (or in the current directory outside git), with the folders, a `SPEC.md` brief stub (Context, Audience, Principles, Design, Stack, Ask first), `notes/glossary.md` and `ROADMAP.md`. It adds `/spec/` to `.gitignore` unless git already ignores it. Existing files are kept, so running it again is safe.

## new

```sh
xspec new <type> <title> [--bug] [--with 7,9] [--dir]
```

Creates an item from its template and prints its path.

| Type | Creates |
| --- | --- |
| `change` | `changes/NNN-slug.md` |
| `update` | `updates/NNN-slug.md`, or a folder with `update.md` + `design.md` with `--dir` |
| `note` | `notes/slug.md` (no number, no status) |
| `analysis` | `analysis/slug.md` (no number, no status) |

`feature`, `issue`, `fix` and `plan` (v0.4) are refused with a message naming `change` and `update`.

- `--bug`: change only. Current / Expected / Unchanged behavior instead of Goal.
- `--with 7,9`: update only. Active changes to include. Run in a terminal without `--with`, `new update` shows a numbered list and asks which to include.
- `--dir`: update only. Makes a folder with `design.md`.

## list

```sh
xspec list [--all]
```

Prints active numbered items: number, type, status, the count of open `[?: …]` markers, checked/total task boxes, and title. `--all` adds archived items.

## next

```sh
xspec next
```

Prints one line: the next unchecked task of the lowest item in progress (an open item with a checked task); else an item whose tasks are all checked, to close; else an item waiting for answers; else the first task of an item not yet started; else `nothing to do`.

## questions

```sh
xspec questions
```

Prints every open `[?: …]` marker in active items as `file:line  question`, grouped by item. Markers inside backticks or code blocks are examples and do not count.

## done

```sh
xspec done <num|slug> [--change <commit|PR|path>] [--decision "text"] [--no-verify]
```

For the item and every active item it includes, it:

1. refuses while a `[?: …]` marker or an unchecked task box remains;
2. runs every `Verify: <command>` line once, in the project root, and refuses if one fails (`--no-verify` skips them);
3. appends `- NNN title (date): text` to `spec/notes/decisions.md` when `--decision` is given, creating the file;
4. ticks the Includes boxes;
5. sets `status: done` and writes `Change: <commit> <YYYY-MM-DD>` on line 2 (the commit from `--change`, else the short git HEAD, else only the date);
6. moves the file or folder to `spec/archive/<type folder>/`.

## roadmap

```sh
xspec roadmap
```

Rewrites and prints `spec/ROADMAP.md`. Every `new` and `done` rewrites it too. Sections: Now (open items with a checked task), Next (the other open items, with question counts), History (archived, newest first, with date and commit). Items included in an update or fix are listed under it.

## context

```sh
xspec context <num|slug>
```

Prints everything an agent needs for one item: the item (and `design.md`), the items it includes, the Principles, Ask first and Stack sections of `SPEC.md`, and `notes/glossary.md`, each under a `=== path ===` header. Paste it into an agent that cannot read `spec/`.

## check

```sh
xspec check
```

Prints every problem and exits 1 when there is one: `spec/` not ignored by git, a bad status line (only `open` and `done` are valid), a duplicate number, work begun (a checked task) with open markers, `done` outside the archive or an archived item that is not done, an Includes line pointing at a missing number, and a requirement (`- R1: …`) that no task box references. It fits a pre-commit hook.

## install

```sh
xspec install project|global [--force]
```

Writes the bundled execspec skill into `.claude/skills/execspec/` and `.agents/skills/execspec/`, under the project's git top-level (`project`) or your home directory (`global`). It refuses to touch an existing skill directory unless `--force` is set; `--force` replaces it whole.
