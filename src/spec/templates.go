package spec

// Placeholders: {{num}}, {{title}}, {{includes}}, {{bugs}}.
// Every [?: …] marker blocks `xspec start` until someone answers it.

const specStub = `# %s

> The owner fills this file. Agents read it and propose edits in chat.

## Goal

## Scope

## Principles

One line each. Plans check themselves against these in their Gates section.

## Decisions
`

const glossaryStub = `# Glossary

Project slang, one line per term: **term** - meaning (where it came up)
`

var templates = map[string]string{
	"feature": `status: open
# {{num}} {{title}}

## Why
[?: the problem, and who gets what when it is solved]

## Requirements
- R1: WHEN [?: condition] THE SYSTEM SHALL [?: behavior]

## Design
[?: the approach and the code it touches]

## Tasks
- [ ] T1 [?: task] (R1)
Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
`,

	"issue": `status: open
# {{num}} {{title}}

## Current behavior
- WHEN [?: condition] THEN the system [?: wrong behavior]

## Expected behavior
- WHEN [?: condition] THEN the system SHALL [?: correct behavior]

## Unchanged behavior
- WHEN [?: condition] THEN the system SHALL CONTINUE TO [?: behavior that must not break]

## Cause
[?: root cause, once found]

## Tasks
- [ ] T1 regression test reproduces the bug
- [ ] T2 [?: fix]
Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
`,

	"update": `status: open
# {{num}} {{title}}

## Goal
[?: why this batch, and what the user gets]

## Includes
{{includes}}
## Parts

### P1 [?: part name]
[?: why this part]
- R1: WHEN [?: condition] THE SYSTEM SHALL [?: behavior]

## Also changes
Migrations, config, dependencies outside the parts.

## Tasks
- [ ] T1 [?: task] (R1)

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]

## Risks
What may break, and how to roll back.
`,

	"fix": `status: open
# {{num}} {{title}}

## Includes
{{includes}}
## Small bugs
{{bugs}}
## Tasks
- [ ] T1 a regression test for every bug above
Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
`,

	"plan": `status: open
# {{num}} {{title}}

## Context
[?: what prompted this, and the constraints]

## Gates
One line per principle in SPEC.md: kept, or broken and why.

## Decision
[?: the chosen approach]

## Steps
- [ ] [?: step; name the items it creates, e.g. 007]
`,

	"feature/light": `status: open
# {{num}} {{title}}

## Goal
[?: what changes, and for whom]

## Tasks
- [ ] T1 [?: task]
`,

	"issue/light": `status: open
# {{num}} {{title}}

## Current behavior
- WHEN [?: condition] THEN the system [?: wrong behavior]

## Expected behavior
- WHEN [?: condition] THEN the system SHALL [?: correct behavior]

## Tasks
- [ ] T1 [?: fix]
`,

	"plan/light": `status: open
# {{num}} {{title}}

## Context
[?: what prompted this]

## Decision
[?: the chosen approach]

## Steps
- [ ] [?: step]
`,

	// Written by init into plans/001-kickoff.md.
	"kickoff": `status: open
# 001 Kickoff

In a project that already has code, describe the current state as the baseline.

## Vision
[?: what it is, for whom, and why now]

## MVP
- R1: WHEN [?: condition] THE SYSTEM SHALL [?: behavior]

## Out of scope
[?: what this project will not do]

## Stack
| Part | Choice | Owner knows it |
| --- | --- | --- |
| [?: language] | [?: choice] | [?: yes or no] |

## Architecture sketch
[?: the main parts and how they connect]

## Milestones
- v0.1: [?: the first usable version]

## Tasks
- [ ] T1 [?: the project skeleton builds and runs] (R1)
`,

	"analysis": `# {{title}}

## Question

## Sources

## Findings

## Verdict
`,

	"note": `# {{title}}
`,

	// design.md beside update.md, and fix.md with --dir.
	"design": `# Design: {{title}}
`,
}
