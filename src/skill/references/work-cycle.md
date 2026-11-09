# Work cycle

For work that does not fit in one sentence: a feature, a bug of average size or more, a batch of work.

## Items

Every item has a number from one counter shared by both types (`007`), so `xspec done 7` and `--with 7,9` work. Names are `<NNN>-<slug>`. The owner may write any item by hand; `xspec new` still gives the number and status line, and `xspec check` verifies the result.

| Type | Where | Holds |
| --- | --- | --- |
| change | `changes/007-login-form.md` | Goal / Tasks / Verify |
| change `--bug` | `changes/009-crash-on-save.md` | Current / Expected / Unchanged behavior / Tasks / Verify |
| update | `updates/012-v0-3.md`; a folder with `update.md` + `design.md` with `--dir` | Goal / Includes / Parts / Tasks / Verify |

## Choosing the shape

- **One piece of work** is one change. A bug is `--bug`.
- **Several pieces given at once** are one update. Each piece is a part: `### P1 <name>`, a line on why, then its tasks. Small bugs are parts too. A piece gets its own change only when it is large, needs its own design, or already exists; list it under Includes.
- **The form grows.** Start with Goal, Tasks, Verify. Add Requirements (EARS) when behavior needs testable wording, Design when the approach is not obvious, Risks for a migration. Add a section only when it changes what you do.

## Format

- Line 1 is the status: `open` or `done`. Line 2 is `Change: <commit> <YYYY-MM-DD>` once done. Then `# 007 Title`. An open item with a checked task is in progress; no status changes.
- Requirements, when present, use EARS form, numbered across the item: `- R1: WHEN <condition> THE SYSTEM SHALL <behavior>`. Every requirement needs a task that names it: `- [ ] T1 add the form (R1)`.
- Bug behavior: `WHEN … THEN the system …` (current), `… SHALL …` (expected), `… SHALL CONTINUE TO …` (unchanged: what the fix must not break).
- Tasks are checkboxes: `- [ ] T1 add the form`.
- An open question is a marker in place: `[?: which auth method?]`. A marker inside backticks or a code block is an example, not a question.
- `Verify: <command>` lines run before `done` closes the item; a failure blocks it. Delete the line when nothing fits.
- Includes lines link changes by number: `- [ ] 007 change: Login form`. Included changes stay in `changes/`.

## Unfamiliar stack

`SPEC.md`'s Stack table says which parts the owner knows. For every choice touching a part they do not know:

- explain it in plain words: what it is, why here, what changes with the alternative;
- leave the decision to the owner as a marker with options and your pick: `[?: SQLite or a JSON file? I recommend SQLite: …]`;
- add each new term to `notes/glossary.md`.

## Steps

1. **Draft.** Read the code the item touches, then `xspec new change "<title>"` (`--bug` for a bug) or `xspec new update "<title>"` (`--with 7,9` to include changes, `--dir` for a design file). Fill the sections; answer the markers the code answers; leave a marker for every question only the owner can answer. If the owner asks to improve it, follow [polish.md](polish.md) first.
2. **Digest.** Show the digest in chat (format in [polish.md](polish.md)), then start. Do not wait for approval.
3. **Build.** Write the code; tick each task box as it is done. `xspec next` names the next unchecked task. Stop and ask the owner in chat at an open marker or an Ask first area; write the answer into the file and remove the marker.
4. **Close.** `xspec done <num>` refuses while a box is unchecked, a marker remains or a Verify command fails. It writes `Change: <HEAD> <date>` (or `--change`) and moves the item, with everything it includes, into `archive/`. Add `--decision "<what was chosen and why>"` when the choice would surprise a later reader; it lands in `notes/decisions.md`.

`spec/ROADMAP.md` is rewritten after every step: Now, Next, History. `xspec check` lists problems: a bad status line, duplicate numbers, work begun with open markers, `done` outside the archive, broken Includes, a requirement with no task, `spec/` not ignored.

Done when every task is checked and the item sits in `archive/` with `status: done`.
