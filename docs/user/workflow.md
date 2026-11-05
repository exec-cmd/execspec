# Workflow

execspec splits a project's knowledge in two. The code and `docs/` hold the *as-built* state: public, tracked, reviewed. `spec/` holds the *intent*: what the project is for, what is planned, what is still undecided. That folder is private to you and your agents and stays out of git.

## Layout

```
spec/
  SPEC.md          yours: Goal, Scope, Principles, Decisions
  personal/        yours: opinions agents must respect
  analysis/        research on tools, repositories, options
  notes/           notes; glossary.md holds project slang
  features/        001-login-form.md
  issues/          002-crash-on-save.md
  updates/         003-v0-2/update.md + design.md
  fixes/           004-typos.md
  plans/           005-auth.md
  archive/         closed items, same subfolders
```

Agents read `SPEC.md` and `personal/` but only propose changes to them in chat. They keep the other folders themselves.

## Item types

| Type | Use it for | Sections |
| --- | --- | --- |
| feature | a new capability | Why / Requirements / Design / Tasks |
| issue | a large bug or problem | Current / Expected / Unchanged behavior / Cause / Tasks |
| update | a batch of features and issues, such as a release | Goal / Includes / Also changes / Tasks / Risks, plus `design.md` |
| fix | a batch of issues and small bugs | Includes / Small bugs / Tasks |
| plan | a design before big work | Context / Gates / Decision / Steps |

Small changes need no file at all.

## Format

```markdown
status: doing
Change: 1389b69                 ← written by `xspec done`
# 007 Login form

## Requirements
- R1: WHEN the password is wrong 5 times THE SYSTEM SHALL lock the account for 15 minutes

## Tasks
- [x] T1 count failed attempts per account (R1)
- [ ] T2 [?: lock in the database or in the cache?] (R1)
```

- **Status** is always line 1.
- **Requirements** use [EARS](https://alistairmavin.com/ears/) and are numbered `R1, R2…`. Every requirement needs a task that names it, and `xspec check` enforces that.
- **Bug behavior** is written as three lists: what happens now, what should happen, and what must keep working (`SHALL CONTINUE TO`). The last list is the regression guard.
- **Open questions** are `[?: …]` markers placed where the answer belongs. Templates start full of them; `start` refuses until they are gone.
- **Includes** reference items by number (`- [ ] 007 feature: Login form`). Included items stay in their own folders until the batch closes.
- **Gates** in a plan take one line per principle in `SPEC.md`: kept, or broken and why.

## Cycle

1. **Draft.** `xspec new feature "Login form"`. The agent reads the code it touches, fills the sections and leaves a marker for every question only you can answer.
2. **Approve.** You answer the markers and approve the file.
3. **Build.** `xspec start 7`. The agent writes the code and ticks each task as it is done.
4. **Close.** `xspec done 7` records the commit and archives the file.

An update or fix runs the same cycle for everything it includes at once.
