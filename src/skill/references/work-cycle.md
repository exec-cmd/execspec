# Work cycle

For a large feature, a bug above average in size or risk, a batch of work, or a design before big work.

## Items

Every item has a number from one counter shared by all types (`007`), so `xspec start 7` and `--with 7,9` work. Names are `<NNN>-<slug>`.

| Type | Where | Holds |
| --- | --- | --- |
| feature | `features/007-login-form.md` | Why / Requirements / Design / Tasks |
| issue | `issues/009-crash-on-save.md` | Current / Expected / Unchanged behavior / Cause / Tasks |
| update | `updates/012-v0-2/update.md` + `design.md` | Goal / Includes / Also changes / Tasks / Risks |
| fix | `fixes/014-typos.md`, or a folder with `design.md` (`--dir`) | Includes / Small bugs / Tasks |
| plan | `plans/015-auth.md` | Context / Gates / Decision / Steps |

Format:

- Line 1 is `status: open`, `status: doing` or `status: done`. Line 2 is `Change: <commit>` once done. Then `# 007 Title`.
- Requirements use EARS form, numbered: `- R1: WHEN <condition> THE SYSTEM SHALL <behavior>`.
- Bug behavior: `WHEN … THEN the system …` (current), `… SHALL …` (expected), `… SHALL CONTINUE TO …` (unchanged: what the fix must not break).
- Tasks are checkboxes that name the requirements they serve: `- [ ] T1 add the form (R1)`. Every requirement needs a task.
- An open question is a marker in place: `[?: which auth method?]`. Templates start full of them.
- Includes lines link items by number: `- [ ] 007 feature: Login form`. Included items stay in their own folders.
- A plan's Gates section holds one line per principle in `SPEC.md`: kept, or broken and why.

## Steps

1. **Draft.** Read the code the item touches, then `xspec new <type> "<title>"` (update/fix: `--with 7,9`; fix: `--bug "text"`). Fill the sections and replace the markers you can answer. Leave a marker for every question only the owner can answer.
2. **Approve.** Show the file and wait for the owner's ok. Open markers block the ok, and `xspec start` refuses while any remain, in the item or in anything it includes.
3. **Build.** `xspec start <num>` sets `doing` on the item and everything it includes. Write the code; tick each task box as it is done.
4. **Close.** `xspec done <num>` refuses while a box is unchecked. It sets `done`, writes `Change: <HEAD>` (or `--change`), and moves the item into `archive/`. An update or fix closes and archives everything it includes along with it.

`xspec check` lists problems: a bad status line, duplicate numbers, `doing` with open markers, `done` outside the archive, broken Includes, a requirement with no task, `spec/` not ignored.

Done when every task is checked and the item sits in `archive/` with `status: done`.
