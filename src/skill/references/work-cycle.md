# Work cycle

For the kickoff, a large feature, a bug above average in size or risk, a batch of work, or a design before big work.

## Items

Every item has a number from one counter shared by all types (`007`), so `xspec start 7` and `--with 7,9` work. Names are `<NNN>-<slug>`. The owner may write any item by hand; `xspec new` still gives the number and status line, and `xspec check` verifies the result.

| Type | Where | Holds |
| --- | --- | --- |
| kickoff | `plans/001-kickoff.md`, made by `xspec init` | Vision / MVP / Out of scope / Stack / Architecture sketch / Milestones / Tasks |
| feature | `features/007-login-form.md` | Why / Requirements / Design / Tasks; light: Goal / Tasks |
| issue | `issues/009-crash-on-save.md` | Current / Expected / Unchanged behavior / Cause / Tasks; light: Current / Expected / Tasks |
| update | `updates/012-v0-3/update.md` + `design.md` | Goal / Includes / Parts / Also changes / Tasks / Risks |
| fix | `fixes/014-typos.md`, or a folder with `design.md` (`--dir`) | Includes / Small bugs / Tasks |
| plan | `plans/015-auth.md` | Context / Gates / Decision / Steps; light: Context / Decision / Steps |

## Choosing the shape

- **One piece of work** is one feature, issue or plan.
- **Several pieces given at once** are one update. Each piece is a part: `### P1 <name>`, a line on why, then its own numbered requirements. Requirements are numbered across the whole update, so tasks can name any of them. A piece gets its own feature only when it needs its own approval or design, or already exists; link it under Includes. Several bugs at once are one fix: issues under Includes, small ones under Small bugs.
- **Depth.** `--depth light` for a familiar stack, about a day of work and low risk; full (the default) otherwise. The owner's word overrides. Updates and fixes are always full.

## Format

- Line 1 is the status: `open`, `approved <hash>`, `doing`, `done`. Line 2 is `Change: <commit> <YYYY-MM-DD>` once done. Then `# 007 Title`.
- Requirements use EARS form, numbered: `- R1: WHEN <condition> THE SYSTEM SHALL <behavior>`.
- Bug behavior: `WHEN … THEN the system …` (current), `… SHALL …` (expected), `… SHALL CONTINUE TO …` (unchanged: what the fix must not break).
- Tasks are checkboxes that name the requirements they serve: `- [ ] T1 add the form (R1)`. Every requirement needs a task.
- An open question is a marker in place: `[?: which auth method?]`. Templates start full of them. A marker inside backticks or a code block is an example, not a question.
- `Verify: <command>` lines run before `done` closes the item; a failure blocks it. Full templates ask for one; delete the line when nothing fits.
- Includes lines link items by number: `- [ ] 007 feature: Login form`. Included items stay in their own folders.
- A plan's Gates section holds one line per principle in `SPEC.md`: kept, or broken and why.

## Unfamiliar stack

The kickoff's Stack table says which parts the owner knows. For every choice touching a part they do not know:

- explain it in plain words: what it is, why here, what changes with the alternative;
- leave the decision to the owner as a marker with options and your pick: `[?: SQLite or a JSON file? I recommend SQLite: …]`;
- add each new term to `notes/glossary.md`.

## Steps

1. **Draft.** Read the code the item touches, then `xspec new <type> "<title>"` (update/fix: `--with 7,9`; fix: `--bug "text"`; `--depth light` when it fits). Fill the sections; answer the markers the code answers; leave a marker for every question only the owner can answer. In a project with code, the kickoff describes the current state as the baseline. If the owner asks to improve it, follow [polish.md](polish.md) first. Then show the digest.
2. **Approve.** Show the digest first, then the file or `xspec context <num>`, and the output of `xspec questions`. The owner answers the markers and runs `xspec approve <num>`; you never run it. An update or fix approves everything it includes. Editing an approved item (ticking boxes aside) voids the approval.
3. **Build.** `xspec start <num>` sets `doing` on the item and everything it includes. Write the code; tick each task box as it is done. `xspec next` names the next unchecked task.
4. **Close.** `xspec done <num>` refuses while a box is unchecked or a Verify command fails. It writes `Change: <HEAD> <date>` (or `--change`), and moves the item, with everything it includes, into `archive/`.

`spec/ROADMAP.md` is rewritten after every step: Now, Next, History. `xspec check` lists problems: a bad status line, duplicate numbers, an edit after approval, `doing` with open markers, `done` outside the archive, broken Includes, a requirement with no task, `spec/` not ignored.

Done when every task is checked and the item sits in `archive/` with `status: done`.
