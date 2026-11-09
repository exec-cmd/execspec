# Workflow

execspec splits a project's knowledge in two. The code and `docs/` hold the *as-built* state: public, tracked, reviewed. `spec/` holds the *intent*: what the project is for, what is planned, what is still undecided. That folder is private to you and your agents and stays out of git.

## Layout

```
spec/
  SPEC.md          project brief: Context, Audience, Principles, Design, Stack, Ask first
  ROADMAP.md       generated: Now, Next, History
  personal/        yours: opinions agents must respect
  analysis/        research on tools, repositories, options
  notes/           notes; glossary.md holds project slang, decisions.md the reasons behind closed items
  changes/         001-login-form.md
  updates/         002-v0-2.md, or 002-v0-2/update.md + design.md
  archive/         closed items, same subfolders
```

Agents read `SPEC.md` and `personal/` but write to them only when you say yes. They keep the other folders themselves.

## The brief

`SPEC.md` has six sections. You do not fill a blank page: on first use the agent drafts them from the README, `AGENTS.md`, the manifest and the code, asks you three or four questions it cannot answer from the repository, shows the draft in chat and writes it when you say yes.

| Section | Holds |
| --- | --- |
| Context | what the project is and why it exists |
| Audience | who uses it |
| Principles | the rules every change follows, one line each |
| Design | a few lines on how it is built, or a link to the design doc |
| Stack | a table of part, choice and whether you know it |
| Ask first | what the agent must stop and ask about; the default: add a dependency, change a file format or public interface, migrate or delete data, touch secrets, CI or release configuration |

When the agent works on a stack part you marked "no", it explains each choice in plain words, leaves the decision to you as a marker with options and its recommendation, and adds new terms to the glossary.

## Item types

| Type | Use it for | Sections |
| --- | --- | --- |
| change | one piece of work: a feature or a bug | Goal / Tasks / Verify |
| change `--bug` | a bug that needs a regression guard | Current / Expected / Unchanged behavior / Tasks / Verify |
| update | several pieces of work at once, such as a release | Goal / Includes / Parts / Tasks / Verify; `--dir` adds `design.md` |

A task you can describe in one sentence needs no file at all. The form grows with the work: add Requirements, Design or Risks when they change what happens next, not by default.

### One update, not a pile of changes

When you hand over several things at once ("tags, due dates, export, and fix the sort bug"), they go into one update. Each piece is a part:

```markdown
## Parts

### P1 Tags
Group tasks by topic.

### P2 Due dates
List overdue tasks first.
```

One file, one place to see how the parts connect. A piece gets its own change only when it is large or already exists; then it goes under Includes.

## Format

```markdown
status: open
# 007 Login form

## Goal
Lock an account after repeated wrong passwords.

## Requirements
- R1: WHEN the password is wrong 5 times THE SYSTEM SHALL lock the account for 15 minutes

## Tasks
- [x] T1 count failed attempts per account (R1)
- [ ] T2 lock in the database (R1)

Verify: go test ./internal/auth/...
```

- **Status** is always line 1: `open` or `done`. A task checked on an open item means it is in progress. Once done, line 2 is `Change: <commit> <date>`.
- **Requirements** are optional. When present they use [EARS](https://alistairmavin.com/ears/) and are numbered `R1, R2…`; every requirement needs a task that names it, and `xspec check` enforces that.
- **Bug behavior** is three lists: what happens now, what should happen, and what must keep working (`SHALL CONTINUE TO`). The last list is the regression guard.
- **Open questions** are `[?: …]` markers placed where the answer belongs. `done` refuses while any remains. Markers inside backticks or code blocks are examples.
- **Verify** lines are commands that must pass before `done` closes the item.
- **By hand.** Items are plain Markdown. Create one with `xspec new` for the number and status line, edit it, then run `xspec check`.
- **Includes** reference changes by number (`- [ ] 007 change: Login form`). Included changes stay in `changes/` until the update closes.

## Cycle

1. **Draft.** The agent runs `xspec new …`, reads the code it touches, fills the sections and leaves a marker for every question only you can answer. `xspec questions` lists them.
2. **Digest.** The agent says in chat, in at most 10 lines, what it will do, and starts. There is nothing to approve.
3. **Build.** The agent writes the code and ticks each task. It stops and asks you when it meets an open marker or work in an Ask first area; the answer goes into the file and the marker goes away. `xspec next` always names the next step.
4. **Close.** `xspec done 7` runs the Verify commands, records the commit and date, and archives the file. `--decision "text"` adds one line of reasoning to `notes/decisions.md`.

An update closes everything it includes at once. `spec/ROADMAP.md` follows every step.

You can ask the agent to improve or summarize any spec file. It changes `SPEC.md` only when you say yes.

## Finding work

Ask the agent to **find bugs** and it reviews the code, writes confirmed large bugs as `change --bug` and small ones as parts of one update, and keeps unconfirmed leads in an analysis. Ask **what to do next** and it writes an analysis of candidates (value, cost, risk) with a recommended set; on your pick it creates one update.
