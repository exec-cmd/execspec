# Workflow

execspec splits a project's knowledge in two. The code and `docs/` hold the *as-built* state: public, tracked, reviewed. `spec/` holds the *intent*: what the project is for, what is planned, what is still undecided. That folder is private to you and your agents and stays out of git.

## Layout

```
spec/
  SPEC.md          yours: Goal, Scope, Principles, Decisions
  ROADMAP.md       generated: Now, Next, History
  personal/        yours: opinions agents must respect
  analysis/        research on tools, repositories, options
  notes/           notes; glossary.md holds project slang
  plans/           001-kickoff.md, 005-auth.md
  features/        002-login-form.md
  issues/          003-crash-on-save.md
  updates/         004-v0-2/update.md + design.md
  fixes/           006-typos.md
  archive/         closed items, same subfolders
```

Agents read `SPEC.md` and `personal/` but only propose changes to them in chat. They keep the other folders themselves.

## Item types

| Type | Use it for | Sections |
| --- | --- | --- |
| kickoff | the start of the project; made by `init` | Vision / MVP / Out of scope / Stack / Architecture sketch / Milestones / Tasks |
| feature | one new capability | Why / Requirements / Design / Tasks |
| issue | one large bug or problem | Current / Expected / Unchanged behavior / Cause / Tasks |
| update | several pieces of work at once, such as a release | Goal / Includes / Parts / Also changes / Tasks / Risks, plus `design.md` |
| fix | several bugs at once | Includes / Small bugs / Tasks |
| plan | a design before big work | Context / Gates / Decision / Steps |

Small changes need no file at all. Features, issues and plans also come in a light depth (`--depth light`) for small, low-risk work in a stack you know.

### One update, not a pile of features

When you hand over several things at once ("tags, due dates, export, and fix the sort bug"), they go into one update. Each piece is a part with its own requirements:

```markdown
## Parts

### P1 Tags
Group tasks by topic.
- R1: WHEN the user runs `todo add "x" --tag work` THE SYSTEM SHALL store the tag

### P2 Due dates
- R2: WHEN a task is past its due date THE SYSTEM SHALL list it first
```

One file, one approval, one place to see how the parts connect. A piece gets its own feature only when it needs its own approval or already exists; then it goes under Includes.

## Format

```markdown
status: approved 3f9a1c2e      ← written by `xspec approve`
# 007 Login form

## Requirements
- R1: WHEN the password is wrong 5 times THE SYSTEM SHALL lock the account for 15 minutes

## Tasks
- [x] T1 count failed attempts per account (R1)
- [ ] T2 lock in the database (R1)

Verify: go test ./internal/auth/...
```

- **Status** is always line 1: `open → approved → doing → done`. Once done, line 2 is `Change: <commit> <date>`.
- **Requirements** use [EARS](https://alistairmavin.com/ears/) and are numbered `R1, R2…`. Every requirement needs a task that names it, and `xspec check` enforces that.
- **Bug behavior** is written as three lists: what happens now, what should happen, and what must keep working (`SHALL CONTINUE TO`). The last list is the regression guard.
- **Open questions** are `[?: …]` markers placed where the answer belongs. Templates start full of them; `approve` refuses until they are gone. Markers inside backticks or code blocks are examples.
- **Verify** lines are commands that must pass before `done` closes the item.
- **Includes** reference items by number (`- [ ] 007 feature: Login form`). Included items stay in their own folders until the batch closes.
- **Gates** in a plan take one line per principle in `SPEC.md`: kept, or broken and why.

## The kickoff and an unfamiliar stack

`xspec init` creates `plans/001-kickoff.md`: what the project is, the MVP, what is out of scope, the stack, milestones. Its Stack table has a column "Owner knows it". When the agent works on a part you marked "no", it explains each choice in plain words, leaves the decision to you as a marker with options and its recommendation, and adds new terms to the glossary. You stay the one who decides, even in a stack you are learning.

In a project that already has code, the kickoff describes the current state as the baseline.

## Cycle

1. **Draft.** The agent runs `xspec new …`, reads the code it touches, fills the sections and leaves a marker for every question only you can answer. `xspec questions` lists them.
2. **Approve.** You answer the markers, read the file (or `xspec context 7`) and run `xspec approve 7`. Only you run it: an edit afterwards voids the approval.
3. **Build.** `xspec start 7`. The agent writes the code and ticks each task as it is done. `xspec next` always names the next step.
4. **Close.** `xspec done 7` runs the Verify commands, records the commit and date, and archives the file.

An update or fix runs the same cycle for everything it includes at once. `spec/ROADMAP.md` follows every step.

## Finding work

Ask the agent to **find bugs** and it reviews the code, writes confirmed large bugs as issues and small ones as one fix, and keeps unconfirmed leads in an analysis. Ask **what to do next** and it writes an analysis of candidates (value, cost, risk) with a recommended set; on your pick it creates one update.
