# Usage

Every command except `init` and `install` looks for `spec/` in the current directory and then in each parent directory.

## init

```sh
xspec init
```

Creates `spec/` at the git top-level (or in the current directory outside git), with the folders, a `SPEC.md` stub and `notes/glossary.md`. It adds `spec/` to `.gitignore` unless git already ignores it. Existing files are kept, so running it again is safe.

## new

```sh
xspec new <type> <title> [--with 7,9] [--bug "text"]... [--dir]
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

- `--with 7,9`: items to include. An update takes features and issues; a fix takes issues. Run in a terminal without `--with`, `new update` and `new fix` show a numbered list and ask which to include.
- `--bug "text"`: fix only, repeatable. Adds a line under Small bugs.
- `--dir`: fix only. Makes a folder with `design.md`.

## list

```sh
xspec list [--all]
```

Prints active numbered items: number, type, status, the count of open `[?: …]` markers, checked/total task boxes, and title. `--all` adds archived items.

## start

```sh
xspec start <num|slug>
```

Sets `status: doing` on the item and on every active item it includes. It refuses while any of them has a `[?: …]` marker.

## done

```sh
xspec done <num|slug> [--change <commit|PR|path>]
```

For the item and every active item it includes, it:

1. refuses while a `[?: …]` marker or an unchecked task box remains;
2. ticks the Includes boxes;
3. sets `status: done` and writes `Change:` on line 2 (from `--change`, or the short git HEAD, or nothing outside git);
4. moves the file or folder to `spec/archive/<type folder>/`.

## check

```sh
xspec check
```

Prints every problem and exits 1 when there is one. It reports: `spec/` not ignored by git, a bad status line, a duplicate number, `doing` with open markers, `done` outside the archive or an archived item that is not done, an Includes line pointing at a missing number, and a requirement (`- R1: …`) that no task box references. It fits a pre-commit hook.

## install

```sh
xspec install project|global [--force]
```

Writes the bundled execspec skill into `.claude/skills/execspec/` and `.agents/skills/execspec/`, under the project's git top-level (`project`) or your home directory (`global`). It refuses to touch an existing skill directory unless `--force` is set; `--force` replaces it whole.
